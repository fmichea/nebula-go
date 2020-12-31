package registers

import (
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

// TDSFlag - Tile Data Select
type TDSFlag struct {
	registerslib.Flag
}

func NewTDSFlag(reg registerslib.Byte, bit uint8) *TDSFlag {
	return &TDSFlag{
		Flag: registerslib.NewFlag(reg, bit),
	}
}

func (f *TDSFlag) AdjustTileAddress(tileNumber uint8) uint8 {
	// Case 1, base is 0x8000 and tileNumber is between 0:255.
	if f.GetBool() {
		return tileNumber
	}

	// Case 2, base is 0x8800 but tileNumber is between -128:127 so 0 is at 0x9000.
	r16 := int16(int8(tileNumber))
	return uint8(128 + r16)
}

func (f *TDSFlag) BaseAddress() uint16 {
	if f.GetBool() {
		return 0x8000
	}
	return 0x8800
}
