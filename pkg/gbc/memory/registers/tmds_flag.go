package registers

import (
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

// TMDSFlag - Tile Map Display Select Flag
type TMDSFlag struct {
	registerslib.Flag
}

func NewTMDSFlag(reg registerslib.Byte, bit uint8) *TMDSFlag {
	return &TMDSFlag{
		Flag: registerslib.NewFlag(reg, bit),
	}
}

func (f *TMDSFlag) BaseAddress() uint16 {
	if f.GetBool() {
		return 0x9C00
	}
	return 0x9800
}
