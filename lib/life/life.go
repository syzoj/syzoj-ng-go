package life

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func SignalContext() context.Context {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-ch
		cancel()
	}()
	return ctx
}
