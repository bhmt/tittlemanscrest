package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/bhmt/tittlemanscrest/api"
	"github.com/bhmt/tittlemanscrest/api/handlers"
	"github.com/bhmt/tittlemanscrest/cmd"
)

func Work(ctx context.Context) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(slog.String("app_id", "example"))

	handler := api.MiddlewareBase(
		logger,
		api.MiddlewareRest(
			http.HandlerFunc(handlers.Health()),
		),
	)

	mux := http.NewServeMux()
	mux.Handle(
		"/health",
		handler,
	)

	server := api.New(":8081", mux)
	go func() { server.ListenAndServe() }()
	logger.InfoContext(ctx, "listening on :8081")

	<-ctx.Done()
}

func main() {
	cmd.Run(Work)
}
