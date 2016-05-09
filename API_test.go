package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestIndexNoAuth(t *testing.T) {
	recorder := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if (err != nil) {
		panic("ack")
	}

	Index(recorder, r)

	h := recorder.Header().Get("WWW-Authenticate")
	if (h != "Basic realm=\"user\"") {
		t.Errorf("Expected WWW-Authenticate but got %v", h)

	}

	if (recorder.Code != http.StatusUnauthorized) {
		t.Errorf("Expected StatusUnauthorized but got %v", http.StatusText(recorder.Code))
	}
}
func TestIndexBadAuth(t *testing.T) {
	recorder := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if (err != nil) {
		panic("ack")
	}
	r.SetBasicAuth("invalid", "invalid")
	Index(recorder, r)

	if (recorder.Code != http.StatusForbidden) {
		t.Errorf("Expected StatusUnauthorized but got %v", http.StatusText(recorder.Code))
	}
}
func TestIndexGoodAuth(t *testing.T) {
	recorder := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if (err != nil) {
		panic("ack")
	}
	r.SetBasicAuth("admin", "pass")
	Index(recorder, r)

	if (recorder.Code != http.StatusOK) {
		t.Errorf("Expected StatusOK but got %v", http.StatusText(recorder.Code))
	}
}
