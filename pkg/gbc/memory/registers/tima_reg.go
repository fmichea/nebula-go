package registers

import (
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

type TIMAReg struct {
	registerslib.Byte

	div   *DIVReg
	tma   *TMAReg
	ifReg *InterruptReg

	overflowHappened bool
}

func NewTIMAReg(ptr *uint8, div *DIVReg, tma *TMAReg, ifReg *InterruptReg) *TIMAReg {
	return &TIMAReg{
		Byte: registerslib.NewByte(ptr, 0x00),

		div:   div,
		tma:   tma,
		ifReg: ifReg,
	}
}

func (r *TIMAReg) Set(value uint8) {
	r.Byte.Set(value)
	r.overflowHappened = false
}

func (r *TIMAReg) MaybeIncrement() {
	if r.div.RetrieveTIMAIncrementRequest() {
		value := r.Get()
		r.Set(value + 1)

		r.overflowHappened = value == 0xFF
	} else if r.overflowHappened {
		r.ifReg.Timer.SetBool(true)
		r.Set(r.tma.Get())
	}
}
