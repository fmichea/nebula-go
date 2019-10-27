package alu

import (
	"nebula-go/pkg/gbc/memory/registers"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) orConstToA(cst uint8) {
	value := f.regs.A.Get() | cst

	f.regs.A.Set(value)

	f.regs.F.Set(z80lib.FlagsCleared)
	f.regs.F.ZF.SetBool(value == 0)
}

func (f *Factory) OrByteToA(reg registers.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.orConstToA(reg.Get())
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}

func (f *Factory) OrHLPtrToA() opcodeslib.Opcode {
	return f.buildHLPtrToAFunc(f.orConstToA)
}

func (f *Factory) OrD8ToA() opcodeslib.Opcode {
	return f.buildD8ToAFunc(f.orConstToA)
}
