package handlers

import (
	"io"
	"net/http"
)

func ChunkedTransferEncoding(fn func() io.Reader) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Connection", "Keep-Alive")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		reader := fn()
		chunk := make([]byte, 256)
		for {
			n, err := reader.Read(chunk)
			if err != nil {
				break
			}
			_, err = w.Write(chunk[:n])
			if err != nil {
				continue
			}

			flusher.Flush()
		}
	}
}
