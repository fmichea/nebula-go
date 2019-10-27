package controlflow

import (
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/memory/registers"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func returnImpl(mmu memory.MMU, regs *z80lib.Registers) opcodeslib.OpcodeResult {
	value, err := popReturnAddress(mmu, regs)
	if err != nil {
		return opcodeslib.OpcodeError(err)
	}
	regs.PC = value
	return opcodeslib.OpcodeSuccess(1, 16)
}

func returnIf(mmu memory.MMU, regs *z80lib.Registers, cond bool) opcodeslib.OpcodeResult {
	if cond {
		value, err := popReturnAddress(mmu, regs)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}
		regs.PC = value
		return opcodeslib.OpcodeSuccess(1, 20)
	}
	return opcodeslib.OpcodeSuccess(1, 8)
}

func (f *Factory) ReturnIf(flag registers.Flag) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return returnIf(f.mmu, f.regs, flag.GetBool())
	}
}

func (f *Factory) ReturnIfNot(flag registers.Flag) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return returnIf(f.mmu, f.regs, !flag.GetBool())
	}
}

func (f *Factory) Return() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return returnImpl(f.mmu, f.regs)
	}
}

func (f *Factory) ReturnInterrupt() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.regs.IME = true
		return returnImpl(f.mmu, f.regs)
	}
}
