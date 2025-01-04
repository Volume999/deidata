package main

import (
	"bytes"
	"log"
	"testing"
)

func TestGreeting(t *testing.T) {
	// Given
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	// When
	greeting(logger)

	// Then
	expected := "Hello, World!\n"
	got := buf.String()
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}
