package opcodes

import (
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/z80/opcodes/alu"
	"nebula-go/pkg/gbc/z80/opcodes/controlflow"
	"nebula-go/pkg/gbc/z80/opcodes/load"
	"nebula-go/pkg/gbc/z80/opcodes/misc"
	"nebula-go/pkg/gbc/z80/opcodes/misc/cb"
	"nebula-go/pkg/gbc/z80/registers"
)

type Factory struct {
	ALU           *alu.Factory
	CB            *cb.Factory
	ControlFlow   *controlflow.Factory
	Load          *load.Factory
	Miscellaneous *misc.Factory

	mmu  memory.MMU
	regs *registers.Registers
}

func NewFactory(mmu memory.MMU, regs *registers.Registers) *Factory {
	return &Factory{
		ALU:           alu.NewFactory(mmu, regs),
		CB:            cb.NewFactory(mmu, regs),
		ControlFlow:   controlflow.NewFactory(mmu, regs),
		Load:          load.NewFactory(mmu, regs),
		Miscellaneous: misc.NewFactory(mmu, regs),

		mmu:  mmu,
		regs: regs,
	}
}
