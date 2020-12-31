package registers

import (
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

type TMAReg struct {
	registerslib.Byte
}

func NewTMAReg(ptr *uint8) *TMAReg {
	return &TMAReg{
		Byte: registerslib.NewByte(ptr, 0x00),
	}
}
