package main

import (
	"log"
	"os"

	"nebula-go/pkg/gbc/memory"
)

func main() {
	logger := os.Stdout

	if _, err := memory.NewMMU(logger, os.Args[1]); err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
}
