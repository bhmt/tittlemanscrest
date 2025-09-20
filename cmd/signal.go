package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func Run(worker func(context.Context)) {
	ctx, _ /*cancel*/ := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
		os.Kill,
	)
	worker(ctx)
}
