package controlflow

import (
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/memory/registers"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func jumpIf(mmu memory.MMU, regs *z80lib.Registers, cond bool) opcodeslib.OpcodeResult {
	if cond {
		value, err := mmu.ReadDByte(regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		regs.PC = value

		return opcodeslib.OpcodeSuccess(3, 16)
	}

	return opcodeslib.OpcodeSuccess(3, 12)
}

func (f *Factory) JumpIf(flag registers.Flag) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return jumpIf(f.mmu, f.regs, flag.GetBool())
	}
}

func (f *Factory) JumpIfNot(flag registers.Flag) opcodeslib.Opcode {
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
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
