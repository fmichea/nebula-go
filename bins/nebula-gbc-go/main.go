package main

import (
	"io"
	"log"
	"os"

	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/z80"
)

func buildAndRunGB(logger io.Writer, filename string) error {
	mmu, err := memory.NewMMUFromFile(logger, filename)
	if err != nil {
		return err
	}

	cpu := z80.NewCPU(mmu)
	return cpu.Run()
}

func main() {
	if err := buildAndRunGB(os.Stdout, os.Args[1]); err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
}
