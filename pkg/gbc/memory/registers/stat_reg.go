package registers

import (
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

type STATReg struct {
	registerslib.Byte

	Timer *STATTimer

	CoincidenceInterrupt registerslib.Flag // LYC=LY Coincidence Interrupt
	OAMInterrupt         registerslib.Flag
	VBlankInterrupt      registerslib.Flag
	HBlankInterrupt      registerslib.Flag
	CoincidenceFlag      registerslib.Flag
	Mode                 *STATModeBitProxy
}

func NewSTATReg(ptr *uint8) *STATReg {
	reg := registerslib.NewByteWithMask(ptr, 0x85, 0xF8)

	timer := NewSTATTimer()

	return &STATReg{
		Byte: reg,

		Timer: timer,

		CoincidenceInterrupt: registerslib.NewFlag(reg, 6),
		OAMInterrupt:         registerslib.NewFlag(reg, 5),
		VBlankInterrupt:      registerslib.NewFlag(reg, 4),
		HBlankInterrupt:      registerslib.NewFlag(reg, 3),
		CoincidenceFlag:      registerslib.NewFlag(reg, 2),
		Mode:                 NewSTATModeBitProxy(reg, timer),
	}
}
