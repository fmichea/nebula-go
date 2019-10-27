package main

import (
	"io"
	"log"
	"os"

	"nebula-go/pkg/gbc/memory"
)

func buildAndRunGB(logger io.Writer, filename string) error {
	_, err := memory.NewMMUFromFile(logger, filename)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if err := buildAndRunGB(os.Stdout, os.Args[1]); err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
}
