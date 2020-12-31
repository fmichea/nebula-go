package registers

import (
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

type LCDCReg struct {
	registerslib.Byte

	LDE    registerslib.Flag // Bit 7 - LCD Display Enable
	WTMDS  *TMDSFlag         // Bit 6 - Window Tile Map Display Select
	WDE    registerslib.Flag // Bit 5 - Window Display Enable
	BGWTDS *TDSFlag          // Bit 4 - BG & Window Tile Data Select
	BGTMDS *TMDSFlag         // Bit 3 - BG Tile Map Display Select
	OBJSS  *ObjSizeFlag      // Bit 2 - OBJ (Sprite) Size
	OBJSDE registerslib.Flag // Bit 1 - OBJ (Sprite) Display Enable
	BGD    registerslib.Flag // Bit 0 - BG Display
}

func NewLCDCReg(ptr *uint8) *LCDCReg {
	reg := registerslib.NewByte(ptr, 0x91)

	return &LCDCReg{
		Byte: reg,

		LDE:    registerslib.NewFlag(reg, 7),
		WTMDS:  NewTMDSFlag(reg, 6),
		WDE:    registerslib.NewFlag(reg, 5),
		BGWTDS: NewTDSFlag(reg, 4),
		BGTMDS: NewTMDSFlag(reg, 3),
		OBJSS:  NewObjSizeReg(reg, 2),
		OBJSDE: registerslib.NewFlag(reg, 1),
		BGD:    registerslib.NewFlag(reg, 0),
	}
}
