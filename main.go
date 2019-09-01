package main

import (
	"log"
	"os"

	"nebula-go/pkg/gbc/memory"
)

func main() {
	if _, err := memory.NewMMU(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}
