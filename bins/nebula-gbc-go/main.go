package main

import (
	"io"
	"log"
	"os"

	"github.com/pkg/profile"
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/z80"
)

func buildAndRunGB(logger io.Writer, filename string) error {
	profilePath := os.Getenv("PROFILE")

	if profilePath != "" {
		defer profile.Start(profile.ProfilePath(profilePath)).Stop()
	}

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
