package alu

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) orConstToA(cst uint8) {
	value := f.regs.A.Get() | cst

	f.regs.A.Set(value)

	f.regs.F.Set(registers.FlagsCleared)
	f.regs.F.ZF.SetBool(value == 0)
}

func (f *Factory) OrByteToA(reg registerslib.Byte) opcodeslib.Opcode {
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
