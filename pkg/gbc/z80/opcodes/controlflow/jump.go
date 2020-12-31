package controlflow

import (
	"nebula-go/pkg/gbc/memory"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func jumpIf(mmu memory.MMU, regs *registers.Registers, cond bool) opcodeslib.OpcodeResult {
	if cond {
		value, err := mmu.ReadDByte(regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		regs.PC = value

		return opcodeslib.OpcodeSuccess(0, 16) // FIXME: size management
	}

	return opcodeslib.OpcodeSuccess(3, 12)
}

func (f *Factory) JumpIf(flag registerslib.Flag) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return jumpIf(f.mmu, f.regs, flag.GetBool())
	}
}

func (f *Factory) JumpIfNot(flag registerslib.Flag) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return jumpIf(f.mmu, f.regs, !flag.GetBool())
	}
}

func (f *Factory) Jump() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return jumpIf(f.mmu, f.regs, true)
	}
}

func (f *Factory) JumpHL() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.regs.PC = f.regs.HL.Get()
		return opcodeslib.OpcodeSuccess(0, 4) // FIXME: size management
	}
}
