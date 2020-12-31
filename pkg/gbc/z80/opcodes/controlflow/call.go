package controlflow

import (
	"nebula-go/pkg/gbc/memory"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) callInterruptInternal(i z80lib.Interrupt, retAddr uint16) opcodeslib.OpcodeResult {
	// fmt.Printf("Calling interrupt %04Xh\n", i.Addr())
	if err := pushReturnAddress(f.mmu, f.regs, retAddr); err != nil {
		return opcodeslib.OpcodeError(err)
	}
	f.regs.PC = i.Addr()
	return opcodeslib.OpcodeSuccess(0, 16) //FIXME: size management
}

func (f *Factory) CallInterruptInplace(i z80lib.Interrupt) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return f.callInterruptInternal(i, f.regs.PC)
	}
}

// CallInterrupt calls an interrupt. Return address will be pushed on the stack.
func (f *Factory) CallInterrupt(i z80lib.Interrupt) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return f.callInterruptInternal(i, f.regs.PC+1)
	}
}

func callIf(mmu memory.MMU, regs *registers.Registers, cond bool) opcodeslib.OpcodeResult {
	if cond {
		value, err := mmu.ReadDByte(regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		if err := pushReturnAddress(mmu, regs, regs.PC+3); err != nil {
			return opcodeslib.OpcodeError(err)
		}

		regs.PC = value

		return opcodeslib.OpcodeSuccess(0, 24) // FIXME: size management.
	}
	return opcodeslib.OpcodeSuccess(3, 12)
}

// CallIf calls the address after the opcode, if the given flag is true. Return address will be pushed on the stack.
func (f *Factory) CallIf(flag registerslib.Flag) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return callIf(f.mmu, f.regs, flag.GetBool())
	}
}

// CallIfNot calls the address after the opcode, if the given flag is false. Return address will be pushed on the stack.
func (f *Factory) CallIfNot(flag registerslib.Flag) opcodeslib.Opcode {
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
