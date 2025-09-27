package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func ServerSentEvents(fn func() <-chan []byte, liveliness time.Duration) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "sse not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		livelinessTicker := time.NewTicker(liveliness)
		livelinessMsg := []byte(":keepalive\n\n")

		message := fn()

		for {
			select {
			case <-ctx.Done():
				return

			case <-livelinessTicker.C:
				_, err := w.Write(livelinessMsg)
				if err != nil {
					http.Error(w, "err writing message", http.StatusInternalServerError)
					return
				}

				flusher.Flush()

			case msg, ok := <-message:

				if !ok {
					http.Error(w, "err reading message", http.StatusInternalServerError)
					return
				}

				_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
				if err != nil {
					http.Error(w, "err writing message", http.StatusInternalServerError)
					return
				}

				flusher.Flush()
			}
		}
	}
}
