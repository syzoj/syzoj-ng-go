package util

import (
	"context"
)

type ErrorContext struct {
	context.Context
	Error error
	ch    chan struct{}
}

func (ctx *ErrorContext) Err() error {
	return ctx.Error
}

func (ctx *ErrorContext) Done() <-chan struct{} {
	if ctx.ch == nil {
		ctx.ch = make(chan struct{})
		close(ctx.ch)
	}
	return ctx.ch
}
