package alu

import (
	"nebula-go/pkg/gbc/memory/registers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) adcConstToA(cst uint8) {
	cy := uint16(f.regs.F.CY.Get())

	f.logicalOpSetAAndUpdateFlags(cst, false, func(a uint16, value uint16) uint16 {
		return a + value + cy
	})
}

func (f *Factory) AdcByteToA(reg registers.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.adcConstToA(reg.Get())
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}

func (f *Factory) AdcHLPtrToA() opcodeslib.Opcode {
	return f.buildHLPtrToAFunc(f.adcConstToA)

}

func (f *Factory) AdcD8ToA() opcodeslib.Opcode {
	return f.buildD8ToAFunc(f.adcConstToA)
}
