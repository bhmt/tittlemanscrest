package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"github.com/bhmt/tittlemanscrest/api"
	"github.com/bhmt/tittlemanscrest/api/handlers"
	"github.com/bhmt/tittlemanscrest/cmd"
)

func events() <-chan []byte {
	out := make(chan []byte)

	go func() {
		for {
			jitter := rand.IntN(500)
			data := fmt.Sprintf("jitter %dms", jitter)
			out <- []byte(data)
			time.Sleep(2*time.Second + time.Duration(jitter)*time.Millisecond)
		}
	}()

	return out
}

func Work(ctx context.Context) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(slog.String("app_id", "sse"))
	liveliness := 15 * time.Second

	mux := http.NewServeMux()
	mux.Handle(
		"/sse",
		api.MiddlewareBase(
			logger,
			http.HandlerFunc(handlers.ServerSentEvents(events, liveliness)),
		),
	)

	server := api.New(":8081", mux)
	go func() { server.ListenAndServe() }()
	logger.InfoContext(ctx, "listening on :8081")

	<-ctx.Done()
}

func main() {
	cmd.Run(Work)
}
