package api

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"
)

func New(addr string, handler http.Handler, opts ...func(*http.Server)) *http.Server {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	for _, o := range opts {
		o(&s)
	}

	return &s
}

func WithDisableGeneralOptionsHandler(val bool) func(*http.Server) {
	return func(s *http.Server) {
		s.DisableGeneralOptionsHandler = val
	}
}

func WithTLSConfig(cfg *tls.Config) func(*http.Server) {
	return func(s *http.Server) {
		s.TLSConfig = cfg
	}
}

func WithReadTimeout(val time.Duration) func(*http.Server) {
	return func(s *http.Server) {
		s.ReadTimeout = val
	}
}

func WithReadHeaderTimeout(val time.Duration) func(*http.Server) {
	return func(s *http.Server) {
		s.ReadHeaderTimeout = val
	}
}

func WithWriteTimeout(val time.Duration) func(*http.Server) {
	return func(s *http.Server) {
		s.WriteTimeout = val
	}
}

func WithIdleTimeout(val time.Duration) func(*http.Server) {
	return func(s *http.Server) {
		s.IdleTimeout = val
	}
}

func WithMaxHeaderBytes(val int) func(*http.Server) {
	return func(s *http.Server) {
		s.MaxHeaderBytes = val
	}
}

func WithTLSNextProto(val map[string]func(*http.Server, *tls.Conn, http.Handler)) func(*http.Server) {
	return func(s *http.Server) {
		s.TLSNextProto = val
	}
}

func WithConnState(val func(net.Conn, http.ConnState)) func(*http.Server) {
	return func(s *http.Server) {
		s.ConnState = val
	}
}

func WithErrorLog(val *log.Logger) func(*http.Server) {
	return func(s *http.Server) {
		s.ErrorLog = val
	}
}

func WithBaseContext(val func(net.Listener) context.Context) func(*http.Server) {
	return func(s *http.Server) {
		s.BaseContext = val
	}
}

func WithConnContext(val func(ctx context.Context, c net.Conn) context.Context) func(*http.Server) {
	return func(s *http.Server) {
		s.ConnContext = val
	}
}
