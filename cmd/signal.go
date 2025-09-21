package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func Run(worker func(context.Context)) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()
	worker(ctx)
}
