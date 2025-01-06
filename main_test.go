package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHealthcheckDoesNotReturnError(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleHealthcheck)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHandleHealthcheckReturnsExpectedResult(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleHealthcheck)
	handler.ServeHTTP(rr, req)
	want := HealthCheck{AppName: "Deidata", Status: "OK"}
	got := HealthCheck{}
	_ = json.Unmarshal(rr.Body.Bytes(), &got)
	if got != want {
		t.Errorf("handler returned unexpected body: got %v want %v", got, want)
	}
}

func TestNewServer(t *testing.T) {
	s := NewServer()
	if s == nil {
		t.Errorf("NewServer returned an error")
		return
	}
	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	if rr := httptest.NewRecorder(); rr != nil {
		s.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Healthcheck failed for server, status: %v", status)
		}
	}

}
