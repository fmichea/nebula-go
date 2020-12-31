package alu

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) xorConstToA(cst uint8) {
	value := f.regs.A.Get() ^ cst

	f.regs.A.Set(value)

	f.regs.F.Set(0x00)
	f.regs.F.ZF.SetBool(value == 0)
}

func (f *Factory) XorByteToA(reg registerslib.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.xorConstToA(reg.Get())
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}

func (f *Factory) XorHLPtrToA() opcodeslib.Opcode {
	return f.buildHLPtrToAFunc(f.xorConstToA)
}

func (f *Factory) XorD8ToA() opcodeslib.Opcode {
	return f.buildD8ToAFunc(f.xorConstToA)
}
