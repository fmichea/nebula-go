package alu

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) sbcConstToA(value uint8) {
	cy := uint16(f.regs.F.CY.Get())

	f.logicalOpSetAAndUpdateFlags(value, true, func(a uint16, value uint16) uint16 {
		return a - value - cy
	})
}

func (f *Factory) SbcByteToA(reg registerslib.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.sbcConstToA(reg.Get())
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}

func (f *Factory) SbcHLPtrToA() opcodeslib.Opcode {
	return f.buildHLPtrToAFunc(f.sbcConstToA)
}

func (f *Factory) SbcD8ToA() opcodeslib.Opcode {
	return f.buildD8ToAFunc(f.sbcConstToA)
}
