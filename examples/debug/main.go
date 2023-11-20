package main

import (
	"context"
	"log"
	"os"

	"github.com/slack-io/slacker"
)

// Showcasing the ability to toggle the slack Debug option via `WithDebug`

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"), slacker.WithDebug(true))

	definition := &slacker.CommandDefinition{
		Command:     "ping",
		Description: "Ping!",
		Handler: func(ctx *slacker.CommandContext) {
			ctx.Response().Reply("pong")
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
