package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type HealthCheck struct {
	AppName string `json:"appName"`
	Status  string `json:"status"`
}

func handleHealthcheck(w http.ResponseWriter, _ *http.Request) {
	hc := HealthCheck{AppName: "Deidata", Status: "OK"}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(hc)
	if err != nil {
		http.Error(w, "Unexpected Error", http.StatusInternalServerError)
		return
	}
}

type Generator interface {
	Generate(sz int) ([]byte, error)
}

type DataGenerationHandler struct {
	sz  int
	gen Generator
}

func NewDataGenerationHandler(sz int, gen Generator) (*DataGenerationHandler, error) {
	return &DataGenerationHandler{sz: sz, gen: gen}, nil
}

func (h *DataGenerationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.gen == nil {
		http.Error(w, "Generator is not set", http.StatusInternalServerError)
		return
	}
	data, err := h.gen.Generate(h.sz)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error during data generation: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error during data retrieval: %v", err), http.StatusInternalServerError)
		return
	}
}

func NewServer(gen Generator) (http.Handler, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck", handleHealthcheck)
	h, err := NewDataGenerationHandler(100, gen)
	if err != nil {
		return nil, err
	}
	mux.Handle("/data", h)
	return mux, nil
}

// TODO: This should be in some example folder
func run(l *log.Logger) {
	s, err := NewServer(nil)
	if err != nil {
		l.Fatal("Failed to create server: ", err)
		return
	}
	l.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", s)
	if err != nil {
		l.Fatal(err)
	}
}

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	run(logger)
}
