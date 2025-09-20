package api

import "net/http"

type responseWriter struct {
	w          http.ResponseWriter
	StatusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w: w, StatusCode: 200}
}

func (rw responseWriter) Header() http.Header {
	return rw.w.Header()
}

func (rw responseWriter) Write(data []byte) (int, error) {
	return rw.w.Write(data)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
	rw.w.WriteHeader(statusCode)
}
