package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/slack-io/slacker"
)

// Deleting messages via timestamp

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	definition := &slacker.CommandDefinition{
		Command: "ping",
		Handler: func(ctx *slacker.CommandContext) {
			t1, _ := ctx.Response().Reply("about to be deleted")

			time.Sleep(time.Second)

			ctx.Response().Delete(ctx.Event().ChannelID, t1)
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
