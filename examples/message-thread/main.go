package main

import (
	"context"
	"log"
	"os"

	"github.com/slack-io/slacker"
)

// Defining a command with an optional description and example. The handler replies to a thread.

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	definition := &slacker.CommandDefinition{
		Command:     "ping",
		Description: "Ping!",
		Examples:    []string{"ping"},
		Handler: func(ctx *slacker.CommandContext) {
			ctx.Response().Reply("pong", slacker.WithInThread(true))
		},
	}

	bot.AddCommand(definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
