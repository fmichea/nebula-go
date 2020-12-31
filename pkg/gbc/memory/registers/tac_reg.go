package registers

import (
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

type TACRate int

const (
	TACRate4096Hz TACRate = iota
	TACRate262144Hz
	TACRate65536Hz
	TACRate16384Hz
)

func (r TACRate) String() string {
	mapping := map[TACRate]string{
		TACRate4096Hz:   "4096 Hz",
		TACRate16384Hz:  "16384 Hz",
		TACRate65536Hz:  "65536 Hz",
		TACRate262144Hz: "262144 Hz",
	}

	if value, ok := mapping[r]; ok {
		return value
	}
	return "[unknown]"
}

var (
	_tacRateValueToRate = []TACRate{
		TACRate4096Hz,
		TACRate262144Hz,
		TACRate65536Hz,
		TACRate16384Hz,
	}
)

type TACRateBitProxy struct {
	registerslib.BitProxy
}

func NewTACRateBitProxy(reg registerslib.Byte) *TACRateBitProxy {
	return &TACRateBitProxy{
		BitProxy: registerslib.NewBitProxy(reg, 0, 0x3),
	}
}

func (b *TACRateBitProxy) GetRate() TACRate {
	return _tacRateValueToRate[b.Get()]
}

type TACReg struct {
	registerslib.Byte

	Enabled registerslib.Flag
	Rate    *TACRateBitProxy
}

func NewTACReg(ptr *uint8) *TACReg {
	reg := registerslib.NewByte(ptr, 0x00)

	return &TACReg{
		Byte: reg,

		Enabled: registerslib.NewFlag(reg, 2),
		Rate:    NewTACRateBitProxy(reg),
	}
}
