package registers

import (
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

type CGBPaletteIndexReg struct {
	registerslib.Byte

	AutoIncrement registerslib.Flag
	Index         registerslib.BitProxy
}

func NewCGBPaletteIndexReg(ptr *uint8) *CGBPaletteIndexReg {
	reg := registerslib.NewByteWithMask(ptr, 0x00, 0xBF)

	return &CGBPaletteIndexReg{
		Byte: reg,

		AutoIncrement: registerslib.NewFlag(reg, 7),
		Index:         registerslib.NewBitProxy(reg, 0, 0x3F),
	}
}
