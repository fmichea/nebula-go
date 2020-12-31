package registers

import (
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

type ObjSizeFlag struct {
	registerslib.Flag
}

func NewObjSizeReg(reg registerslib.Byte, bit uint8) *ObjSizeFlag {
	return &ObjSizeFlag{
		Flag: registerslib.NewFlag(reg, bit),
	}
}

// return: height in pixels
func (f *ObjSizeFlag) SpriteHeight() int16 {
	if f.GetBool() {
		return 16
	}
	return 8
}
