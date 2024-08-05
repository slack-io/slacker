package slacker

import "github.com/slack-go/slack/socketmode"

// CommandMiddlewareHandler represents the command middleware handler function
type CommandMiddlewareHandler func(CommandHandler) CommandHandler

// CommandHandler represents the command handler function
type CommandHandler func(*CommandContext)

// InteractionMiddlewareHandler represents the interaction middleware handler function
type InteractionMiddlewareHandler func(InteractionHandler) InteractionHandler

// InteractionHandler represents the interaction handler function
type InteractionHandler func(*InteractionContext, *socketmode.Request)

// JobMiddlewareHandler represents the job middleware handler function
type JobMiddlewareHandler func(JobHandler) JobHandler

// JobHandler represents the job handler function
type JobHandler func(*JobContext)
