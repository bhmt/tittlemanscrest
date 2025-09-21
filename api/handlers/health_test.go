package handlers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bhmt/tittlemanscrest/api/handlers"
)

func TestHealth(t *testing.T) {
	health := handlers.Health()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	health(recorder, request)
	data, _ := io.ReadAll(recorder.Result().Body)

	var got map[string]string
	if err := json.Unmarshal(data, &got); err != nil {
		t.Error(err)
	}

	if status := got["status"]; status != "ok" {
		t.Errorf("health data missmatch, want=ok got=%s", status)
	}
}
