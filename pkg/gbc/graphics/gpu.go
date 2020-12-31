package graphics

import (
	"nebula-go/pkg/common/clock"
	"nebula-go/pkg/common/frontends"
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/memory/cartridge"
)

type GPU interface {
	DoCycles(cycles uint16) error
}

type gpu struct {
	mmu     memory.MMU
	mmuRegs *memory.Registers
	cr      *cartridge.ROM

	pacer *Pacer

	display frontends.MainWindow
}

func NewGPU(mmu memory.MMU, display frontends.MainWindow) *gpu {
	return newGPU(mmu, display, NewPacer(clock.New()))
}

func newGPU(mmu memory.MMU, display frontends.MainWindow, pacer *Pacer) *gpu {
	return &gpu{
		mmu:     mmu,
		mmuRegs: mmu.Registers(),
		cr:      mmu.Cartridge(),

		pacer: pacer,

		display: display,
	}
}
