package controlflow

import (
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/memory/registers"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func jumpRelativeIf(mmu memory.MMU, regs *z80lib.Registers, cond bool) opcodeslib.OpcodeResult {
	if cond {
		value, err := mmu.ReadByte(regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		regs.PC = opcodeslib.AddRelativeConst(regs.PC, value)
		return opcodeslib.OpcodeSuccess(2, 12)
	}

	return opcodeslib.OpcodeSuccess(2, 8)
}

func (f *Factory) JumpRelativeIf(flag registers.Flag) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return jumpRelativeIf(f.mmu, f.regs, flag.GetBool())
	}
}

func (f *Factory) JumpRelativeIfNot(flag registers.Flag) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return jumpRelativeIf(f.mmu, f.regs, !flag.GetBool())
	}
}

func (f *Factory) JumpRelative() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return jumpRelativeIf(f.mmu, f.regs, true)
	}
}
