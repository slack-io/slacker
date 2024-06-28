// main.go

package main

import (
	"context"
	"log"
	"os"

	"github.com/slack-io/slacker"
)

// Defining commands using slacker

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.AddCommand(&slacker.CommandDefinition{
		Command: "without",
		Handler: func(ctx *slacker.CommandContext) {
			ctx.Response().Reply("https://www.youtube.com/watch?v=dQw4w9WgXcQ", slacker.WithUnfurlLinks(false))
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
