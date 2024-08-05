package slacker

import "github.com/slack-go/slack/socketmode"

func executeCommand(ctx *CommandContext, handler CommandHandler, middlewares ...CommandMiddlewareHandler) {
	if handler == nil {
		return
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	handler(ctx)
}

func executeInteraction(ctx *InteractionContext, handler InteractionHandler, request *socketmode.Request, middlewares ...InteractionMiddlewareHandler) {
	if handler == nil {
		return
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	handler(ctx, request)
}

func executeJob(ctx *JobContext, handler JobHandler, middlewares ...JobMiddlewareHandler) func() {
	if handler == nil {
		return func() {}
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return func() {
		handler(ctx)
	}
}
