package helper

import (
	"net/http"

	"github.com/google/uuid"
)

var requestIdHeader = "X-Request-Id"

func GetCtxRequestId(r *http.Request) string {
	if id, ok := r.Context().Value(requestIdHeader).(string); ok {
		return id
	}

	return ""
}

func GetHeaderRequestId(r *http.Request) string {
	if id := r.Header.Get(requestIdHeader); id != "" {
		return id
	}

	v7, _ := uuid.NewV7()
	return v7.String()
}
