package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bhmt/tittlemanscrest/handlers"
)

var want = "My name is Jon Daker.\n'Ord is risen today\nAaahhhhggglayloooya\n"
var headers = map[string]string{
	"Connection":             "Keep-Alive",
	"X-Content-Type-Options": "nosniff",
}

func fn() io.Reader {
	return strings.NewReader(want)
}

func TestCTE(t *testing.T) {
	cte := handlers.ChunkedTransferEncoding(fn)

	request := httptest.NewRequest(http.MethodGet, "/cte", nil)
	recorder := httptest.NewRecorder()

	go func() {
		cte(recorder, request)
	}()
	time.Sleep(200 * time.Millisecond)

	result := recorder.Result()

	h_c := result.Header.Get("Connection")
	hg_xcto := result.Header.Get("X-Content-Type-Options")

	if headers["Connection"] != h_c || headers["X-Content-Type-Options"] != hg_xcto {
		t.Errorf("cte haders missmatch, \nwant=%+v\ngot=%+v", headers, result.Header)
	}

	h := headers["Connection"]
	if h != h_c {
		t.Errorf("cte haders missmatch, \nwant=%+v\ngot=%+v", h, h_c)
	}

	output, err := io.ReadAll(result.Body)
	if err != nil {
		t.Error(err)
	}

	got := string(output[:])

	if want != got {
		t.Errorf("cte data missmatch, \nwant=%+v\ngot=%+v", want, got)
	}
}
