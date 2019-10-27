package controlflow

import (
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/memory/registers"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

// CallInterrupt calls an interrupt. Return address will be pushed on the stack.
func (f *Factory) CallInterrupt(i z80lib.Interrupt) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		// TODO: log here about running interrupt?
		if err := pushReturnAddress(f.mmu, f.regs, f.regs.PC+1); err != nil {
			return opcodeslib.OpcodeError(err)
		}
		f.regs.PC = i.Addr()
		return opcodeslib.OpcodeSuccess(1, 16)
	}
}

func callIf(mmu memory.MMU, regs *z80lib.Registers, cond bool) opcodeslib.OpcodeResult {
	if cond {
		value, err := mmu.ReadDByte(regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		if err := pushReturnAddress(mmu, regs, regs.PC+3); err != nil {
			return opcodeslib.OpcodeError(err)
		}

		regs.PC = value

		return opcodeslib.OpcodeSuccess(3, 24)
	}
	return opcodeslib.OpcodeSuccess(3, 12)
}

// CallIf calls the address after the opcode, if the given flag is true. Return address will be pushed on the stack.
func (f *Factory) CallIf(flag registers.Flag) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return callIf(f.mmu, f.regs, flag.GetBool())
	}
}

// CallIfNot calls the address after the opcode, if the given flag is false. Return address will be pushed on the stack.
func (f *Factory) CallIfNot(flag registers.Flag) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return callIf(f.mmu, f.regs, !flag.GetBool())
	}
}

// Call calls the address after the opcode. Return address will be pushed on the stack.
func (f *Factory) Call() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return callIf(f.mmu, f.regs, true)
	}
}
