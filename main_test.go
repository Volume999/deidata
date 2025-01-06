package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO: Test with Dummy Data
type DummyData struct {
	SomeNum  int    `json:"someNum"`
	SomeText string `json:"someText"`
}
type DummyGenerator struct{}

type GenerateFunc func(sz int) ([]byte, error)

func (f GenerateFunc) Generate(sz int) ([]byte, error) {
	return f(sz)
}

func GenerateDummyData(sz int) ([]byte, error) {
	data := make([]DummyData, sz)
	for i := 0; i < sz; i++ {
		data[i] = DummyData{SomeNum: i, SomeText: "SomeText"}
	}
	return json.Marshal(data)
}

func TestDataGenerationHandlerReturnsErrorWhenGeneratorIsNotSet(t *testing.T) {
	req, _ := http.NewRequest("GET", "/data", nil)
	rr := httptest.NewRecorder()
	handler, _ := NewDataGenerationHandler(100, nil)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
	if body := rr.Body.String(); body != "Generator is not set\n" {
		t.Errorf("handler returned unexpected body: got %v want %v", body, "Generator is not set")
	}
}

func TestDataGenerationHandlerReturnsErrorWhenDataGenerationFails(t *testing.T) {
	failGen, _ := NewDataGenerationHandler(
		100,
		GenerateFunc(func(sz int) ([]byte, error) {
			return nil, errors.New("some error")
		}))
	req, _ := http.NewRequest("GET", "/data", nil)
	rr := httptest.NewRecorder()
	failGen.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
	if body := rr.Body.String(); body != "Error during data generation: some error\n" {
		t.Errorf("handler returned unexpected body: got %v want %v", body, "Error during data generation: some error")
	}
}

func TestDataGenerationHandlerReturnsExpectedResult(t *testing.T) {
	gen, _ := NewDataGenerationHandler(2, GenerateFunc(GenerateDummyData))
	req, _ := http.NewRequest("GET", "/data", nil)
	rr := httptest.NewRecorder()
	gen.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	want := []DummyData{{SomeNum: 0, SomeText: "SomeText"}, {SomeNum: 1, SomeText: "SomeText"}}
	var got []DummyData
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Errorf("handler response could not be unmarshalled: %v", err)
	}
	if len(got) != len(want) {
		t.Errorf("handler returned unexpected body: got %v want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("handler returned unexpected body: got %v want %v", got, want)
		}
	}
}

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

func TestNewServerRespondsToHealthcheck(t *testing.T) {
	s, err := NewServer(nil)
	if s == nil {
		t.Errorf("NewServer returned an error: %v", err)
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
