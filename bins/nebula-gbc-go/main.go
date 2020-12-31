package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/profile"

	"nebula-go/pkg/common/frontends"
	"nebula-go/pkg/gbc/graphics"
	graphicslib "nebula-go/pkg/gbc/graphics/lib"
	"nebula-go/pkg/gbc/joypad"
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/z80"
)

func buildAndRunGB(logger io.Writer, filename string) error {
	profilePath := os.Getenv("PROFILE")

	if profilePath != "" {
		defer profile.Start(profile.ProfilePath(profilePath)).Stop()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mmu, err := memory.NewMMUFromFile(logger, filename)
	if err != nil {
		return err
	}

	windowTitle := fmt.Sprintf("nebula-go: %s", mmu.Cartridge().Title)

	window, err := frontends.NewSDLWindow(windowTitle, int32(graphicslib.Width), int32(graphicslib.Height))
	if err != nil {
		return err
	}
	defer func() { _ = window.Close() }()

	gpu := graphics.NewGPU(mmu, window)
	_ = joypad.New(mmu, window)

	cpu := z80.NewCPU(mmu, gpu)
	go func() {
		defer cancel()

		if err := cpu.Run(); err != nil {
			log.Fatalln("CPU error:", err)
		}
	}()

	return window.MainLoop(ctx, cancel)
}

func main() {
	if err := buildAndRunGB(os.Stdout, os.Args[1]); err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
}
