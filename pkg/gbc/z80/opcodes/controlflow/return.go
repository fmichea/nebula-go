package controlflow

import (
	"nebula-go/pkg/gbc/memory"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func returnImpl(mmu memory.MMU, regs *registers.Registers, clock uint16) opcodeslib.OpcodeResult {
	value, err := popReturnAddress(mmu, regs)
	if err != nil {
		return opcodeslib.OpcodeError(err)
	}
	regs.PC = value
	return opcodeslib.OpcodeSuccess(0, clock) // FIXME: size management
}

func returnIf(mmu memory.MMU, regs *registers.Registers, cond bool, clock uint16) opcodeslib.OpcodeResult {
	if cond {
		return returnImpl(mmu, regs, clock)
	}
	return opcodeslib.OpcodeSuccess(1, 8)
}

func (f *Factory) ReturnIf(flag registerslib.Flag) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return returnIf(f.mmu, f.regs, flag.GetBool(), 20)
	}
}

func (f *Factory) ReturnIfNot(flag registerslib.Flag) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return returnIf(f.mmu, f.regs, !flag.GetBool(), 20)
	}
}

func (f *Factory) Return() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return returnImpl(f.mmu, f.regs, 16)
	}
}

func (f *Factory) ReturnInterrupt() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.regs.IME = true
		return returnImpl(f.mmu, f.regs, 16)
	}
}
