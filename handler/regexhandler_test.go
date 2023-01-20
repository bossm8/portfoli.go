package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bossm8.ch/portfolio/handler"
)

func TestRegexHandler(t *testing.T) {
	t.Helper()

	var req *http.Request
	var err error

	//idxHandler := func(w http.ResponseWriter, r *http.Request) {
	//	w.WriteHeader(http.StatusOK)
	//}

	matchAllHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	_http := &handler.RegexHandler{}
	_http.HandleFunc(".*", matchAllHandler)

	if req, err = http.NewRequest(http.MethodGet, "/", nil); err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	_http.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Expected 200, got: %d", status)
	}

}
