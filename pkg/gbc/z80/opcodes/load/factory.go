package load

import (
	"nebula-go/pkg/gbc/memory"
	z80lib "nebula-go/pkg/gbc/z80/lib"
)

type Factory struct {
	mmu  memory.MMU
	regs *z80lib.Registers
}

func NewFactory(mmu memory.MMU, regs *z80lib.Registers) *Factory {
	return &Factory{
		mmu:  mmu,
		regs: regs,
	}
}
