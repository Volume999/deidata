package main

import (
	"encoding/json"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck", handleHealthcheck)
	return mux
}

func run(l *log.Logger) {
	s := NewServer()
	l.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", s)
	if err != nil {
		l.Fatal(err)
	}
}

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	run(logger)
}
