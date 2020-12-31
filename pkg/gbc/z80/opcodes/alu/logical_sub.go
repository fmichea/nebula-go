package alu

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) subConstToA(value uint8) {
	f.logicalOpSetAAndUpdateFlags(value, true, func(a uint16, value uint16) uint16 {
		return a - value
	})
}

func (f *Factory) SubByteToA(reg registerslib.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.subConstToA(reg.Get())
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}

func (f *Factory) SubHLPtrToA() opcodeslib.Opcode {
	return f.buildHLPtrToAFunc(f.subConstToA)
}

func (f *Factory) SubD8ToA() opcodeslib.Opcode {
	return f.buildD8ToAFunc(f.subConstToA)
}
