package api

import (
	"bufio"
	"net"
	"net/http"
)

type intercept struct {
	http.ResponseWriter
	StatusCode int
}

func newIntercept(w http.ResponseWriter) *intercept {
	return &intercept{ResponseWriter: w, StatusCode: 200}
}

func (i intercept) Header() http.Header {
	return i.ResponseWriter.Header()
}

func (i intercept) Write(data []byte) (int, error) {
	return i.ResponseWriter.Write(data)
}

func (i *intercept) WriteHeader(statusCode int) {
	i.StatusCode = statusCode
	i.ResponseWriter.WriteHeader(statusCode)
}

func (i *intercept) Flush() {
	if f, ok := i.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
func (i *intercept) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := i.ResponseWriter.(http.Hijacker); ok {
		h.Hijack()
	}
	return nil, nil, nil
}
