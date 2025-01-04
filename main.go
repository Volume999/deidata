package main

import (
	"log"
	"os"
)

func greeting(l *log.Logger) {
	l.Println("Hello, World!")
}

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	greeting(logger)
}
