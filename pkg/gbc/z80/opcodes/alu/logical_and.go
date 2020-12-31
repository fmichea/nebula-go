package alu

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) andConstToA(cst uint8) {
	value := f.regs.A.Get() & cst

	f.regs.A.Set(value)

	// reset everything, and set HC.
	f.regs.F.Set(registers.HC)
	f.regs.F.ZF.SetBool(value == 0)
}

func (f *Factory) AndByteToA(reg registerslib.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.andConstToA(reg.Get())
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}

func (f *Factory) AndHLPtrToA() opcodeslib.Opcode {
	return f.buildHLPtrToAFunc(f.andConstToA)
}

func (f *Factory) AndD8ToA() opcodeslib.Opcode {
	return f.buildD8ToAFunc(f.andConstToA)
}
