package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bhmt/tittlemanscrest/api/handlers"
)

var sseMsg = "test sse"

func fnSSE() <-chan []byte {
	out := make(chan []byte)

	go func() {
		for {
			out <- []byte(sseMsg)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	return out
}

func TestSSE(t *testing.T) {
	sse := handlers.ServerSentEvents(fnSSE, 15*time.Millisecond)

	request := httptest.NewRequest(http.MethodGet, "/sse", nil)
	recorder := httptest.NewRecorder()

	go func() {
		sse(recorder, request)
	}()
	time.Sleep(100 * time.Millisecond)

	result := recorder.Result()
	data, err := io.ReadAll(result.Body)
	if err != nil {
		t.Error(err)
	}

	var output []string
	for item := range strings.SplitSeq(string(data), "\n\n") {
		if item == "" {
			continue
		}

		output = append(output, item)
	}

	if len(output) == 0 {
		t.Error("sse data missing")
	}

	for _, out := range output {
		got := strings.TrimPrefix(out, "data: ")

		if got == ":keepalive" {
			continue
		}

		if sseMsg != got {
			t.Errorf("sse data missmatch, want=%s\ngot=%s", sseMsg, got)
		}
	}
}
