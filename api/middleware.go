package api

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/bhmt/tittlemanscrest/api/helper"
)

type requestIdContextKeyType struct{}

var requestIdContextKey = requestIdContextKeyType{}

func MiddlewareRest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func MiddlewareBase(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(buf))

		iw := newResponseWriter(w)

		requestId := helper.GetHeaderRequestId(r)
		r = r.WithContext(context.WithValue(r.Context(), requestIdContextKey, requestId))

		start := time.Now()
		logger.LogAttrs(
			r.Context(),
			slog.LevelInfo,
			"request",
			slog.String("request_id", requestId),
			slog.Time("time", start.UTC()),
			slog.String("method", r.Method),
			slog.String("host", r.Host),
			slog.String("path", r.URL.Path),
			slog.String("query", r.URL.RawQuery),
			slog.String("ip", r.RemoteAddr),
		)

		next.ServeHTTP(iw, r)

		end := time.Now()
		logger.LogAttrs(
			r.Context(),
			slog.LevelInfo,
			"response",
			slog.String("request_id", requestId),
			slog.Time("time", end.UTC()),
			slog.Duration("duration", end.Sub(start)),
			slog.Int("status", iw.StatusCode))
	})
}
