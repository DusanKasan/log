package log_test

import (
	"testing"
	"net/http"
	"github.com/DusanKasan/log"
	"strings"
	"io/ioutil"
	"net/http/httptest"
)

func TestMiddleware(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "http://domain.mytld", strings.NewReader("body"))
	if err != nil {
		panic(err)
	}

	response := httptest.NewRecorder()

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("invalid request method in handler. expected: %q, got: %q", http.MethodGet, r.Method)
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		if string(b) != "body" {
			t.Errorf("invalid request body in handler. expected: %q, got %q", "body", string(b))
		}

		if log.GetRequestId(r.Context()) == "" {
			t.Errorf("no request ID ")
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("response"))
	}

	logged := &strings.Builder{}
	log.SetOutput(logged)

	log.Middleware(handler)(response, request)

	if requestID := response.Header().Get("Request-Id"); requestID == "" {
		t.Errorf("request id not set in response header")
	}

	if response.Body.String() != "response" {
		t.Errorf("invalid response body. expected: %q, got %q", "response", response.Body.String())
	}

	if response.Code != http.StatusNotFound {
		t.Errorf("invalid response. expected: %q, got %q", http.StatusNotFound, response.Code)
	}

	logs := strings.Split(strings.Trim(logged.String(), "\n"), "\n")
	if len(logs) != 2 {
		t.Errorf("invalid number of logs. expected: %v, got: %v. details: %v", 2, len(logs), logs)
	}
}