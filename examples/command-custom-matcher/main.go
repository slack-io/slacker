package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/slack-io/commander"
	"github.com/slack-io/proper"
	"github.com/slack-io/slacker"
)

const (
	ignoreCase               = "(?i)"
	wordParameterPattern     = "{\\S+}"
	sentenceParameterPattern = "<\\S+>"
	spacePattern             = "\\s+"
	wordInputPattern         = "(\\S+)"
	sentenceInputPattern     = "(.+)"
	preCommandPattern        = "^"
	postCommandPattern       = "(\\s|$)"
	tokenExtractPattern      = `(?:[^<{]\S+)(?:\s+[^<{]\S+)*`
)

const (
	notParameter      = "NOT_PARAMETER"
	wordParameter     = "WORD_PARAMETER"
	sentenceParameter = "SENTENCE_PARAMETER"
)

func NewCustomCommandMatcher(cmdDef *slacker.CommandDefinition) *CustomCommandMatcher {
	return &CustomCommandMatcher{cmdDef: cmdDef}
}

type CustomCommandMatcher struct {
	cmdDef      *slacker.CommandDefinition
}

func (bc *CustomCommandMatcher) Definition() *slacker.CommandDefinition {
	return bc.cmdDef
}

/*
  Tokenize returns the tokens of the command. This is used to extract the parameters and match the command
  This implementation of tokenize groups multiple words together until a parameter is found
*/
func (bc *CustomCommandMatcher) Tokenize() []*commander.Token {
	tokenExtractRegex := regexp.MustCompile(tokenExtractPattern)
	parameterRegex := regexp.MustCompile(sentenceParameterPattern)
	lazyParameterRegex := regexp.MustCompile(wordParameterPattern)
	cmdDef := bc.cmdDef
	words := tokenExtractRegex.FindAllString(cmdDef.Command, -1)
	tokens := make([]*commander.Token, len(words))
	for i, word := range words {
		word = strings.TrimSpace(word)
		switch {
		case lazyParameterRegex.MatchString(word):
			tokens[i] = &commander.Token{Word: word[1 : len(word)-1], Type: wordParameter}
		case parameterRegex.MatchString(word):
			tokens[i] = &commander.Token{Word: word[1 : len(word)-1], Type: sentenceParameter}
		default:
			tokens[i] = &commander.Token{Word: word, Type: notParameter}
		}
	}
	return tokens
}

/*
  Match returns the parameters of the command and whether or not the command was matched
  NOTE: this doesn't handle @mentions
*/
func (bc *CustomCommandMatcher) Match(text string) (*proper.Properties, bool) {
	tokens := bc.Tokenize()
	expression := bc.generateCommandMatchRegex()

	matches := expression.FindStringSubmatch(text)
	if len(matches) == 0 {
		return nil, false
	}

	if len(matches) != len(tokens) {
		return nil, false
	}

	if len(matches) == 1 {
		return proper.NewProperties(map[string]string{}), true
	}

	parameters := make(map[string]string)
	for i, match := range matches {
		token := tokens[i]
		if !token.IsParameter() {
			continue
		}

		parameters[token.Word] = match
	}
	return proper.NewProperties(parameters), true
}

/*
  generateCommandMatchRegex takes the command definition and returns a regex expression that can be used to
  match the command and extract the parameters. This is called by the Match function to generate the regex
  expression
*/
func (bc *CustomCommandMatcher) generateCommandMatchRegex() *regexp.Regexp {
	expressionString := []string{}
	for _, t := range bc.Tokenize() {
		switch t.Type {
		case wordParameter:
			expressionString = append(expressionString, wordInputPattern)
		case sentenceParameter:
			expressionString = append(expressionString, sentenceInputPattern)
		default:
			expressionString = append(expressionString, regexp.QuoteMeta(t.Word))
		}
	}

	return regexp.MustCompile(ignoreCase + preCommandPattern + strings.Join(expressionString, spacePattern))
}

func hello() *slacker.CommandDefinition{

	return &slacker.CommandDefinition{
		Description: "Echo a message",
		Command:     "hello <message>",
		Examples:    []string{"hello hello"},
		Middlewares: []slacker.CommandMiddlewareHandler{},
		Handler: func(ctx *slacker.CommandContext) {
			message := fmt.Sprintf("Why hello there: %v", ctx.Request().Properties().StringParam("message", "nothing"))
			ctx.Response().Reply(message)
		},
	}
}

func echo() *slacker.CommandDefinition{
	return &slacker.CommandDefinition{
		Description: "Echo a message",
		Command:     "echo <message>",
		Examples:    []string{"hello hello"},
		Middlewares: []slacker.CommandMiddlewareHandler{},
		Handler: func(ctx *slacker.CommandContext) {
			message := fmt.Sprintf("reply: %v", ctx.Request().Properties().StringParam("message", "nothing"))
			ctx.Response().Reply(message)
		},
	}
}

func main() {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")
	server := slacker.NewClient(botToken, appToken)

	personGroup := server.AddCommandGroup("person")
	personGroup.AddCustomCommand(NewCustomCommandMatcher(hello()))
	server.AddCustomCommand(NewCustomCommandMatcher(echo()))
	server.Listen(context.Background())
}