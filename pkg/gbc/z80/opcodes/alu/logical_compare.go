package alu

import (
	"nebula-go/pkg/gbc/memory/registers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) compareConstToA(cst uint8) {
	f.logicalOpUpdateFlags(cst, true, func(a uint16, value uint16) uint16 {
		return a - value
	})
}

func (f *Factory) CompareByteToA(reg registers.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.compareConstToA(reg.Get())
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}

func (f *Factory) CompareHLPtrToA() opcodeslib.Opcode {
	return f.buildHLPtrToAFunc(f.compareConstToA)
}

func (f *Factory) CompareD8ToA() opcodeslib.Opcode {
	return f.buildD8ToAFunc(f.compareConstToA)
}
