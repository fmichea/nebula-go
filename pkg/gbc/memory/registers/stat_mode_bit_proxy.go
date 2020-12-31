package registers

import (
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

type STATModeBitProxy struct {
	registerslib.BitProxy

	timer *STATTimer
}

func NewSTATModeBitProxy(reg registerslib.Byte, timer *STATTimer) *STATModeBitProxy {
	return &STATModeBitProxy{
		BitProxy: registerslib.NewBitProxy(reg, 0, 0x3),

		timer: timer,
	}
}

func (p *STATModeBitProxy) SetMode(mode STATMode) {
	p.timer.SwitchMode(mode)
	p.Set(uint8(mode))
}

func (p *STATModeBitProxy) GetMode() STATMode {
	return STATMode(p.Get())
}
