package load

import (
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/z80/registers"
)

type Factory struct {
	mmu  memory.MMU
	regs *registers.Registers
}

func NewFactory(mmu memory.MMU, regs *registers.Registers) *Factory {
	return &Factory{
		mmu:  mmu,
		regs: regs,
	}
}
