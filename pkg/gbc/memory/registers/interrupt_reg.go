package registers

import (
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

type InterruptReg struct {
	registerslib.Byte

	VBlank *InterruptFlag
	STAT   *InterruptFlag
	Timer  *InterruptFlag
	Serial *InterruptFlag
	Joypad *InterruptFlag
}

func NewInterruptReg(ptr *uint8) *InterruptReg {
	reg := registerslib.NewThreadSafeByte(0x00)

	return &InterruptReg{
		Byte: reg,

		VBlank: NewInterruptFlag(reg, 0),
		STAT:   NewInterruptFlag(reg, 1),
		Timer:  NewInterruptFlag(reg, 2),
		Serial: NewInterruptFlag(reg, 3),
		Joypad: NewInterruptFlag(reg, 4),
	}
}
