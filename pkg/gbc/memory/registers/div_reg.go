package registers

import (
	"nebula-go/pkg/common/bitwise"
)

var (
	// FIXME: document this.
	_divBitForTac = []uint16{9, 3, 5, 7}
)

type DIVReg struct {
	tac              *TACReg
	previousBitValue uint8
	timaIncReq       bool

	cyclesCounter uint16
}

func NewDIVReg(ptr *uint8, tac *TACReg) *DIVReg {
	return &DIVReg{
		timaIncReq: false,

		tac:           tac,
		cyclesCounter: 0,
	}
}

func (r *DIVReg) Get() uint8 {
	return uint8(bitwise.GetBits16(r.cyclesCounter, 8, 0xFF))
}

func (r *DIVReg) SetNoMask(value uint8) {
	r.Set(value)
}

func (r *DIVReg) Set(value uint8) {
	r.cyclesCounter = 0
	r.handleIncrementRequest()
}

func (r *DIVReg) MaybeIncrement(cycles uint16) {
	r.cyclesCounter += cycles
	r.handleIncrementRequest()
}

func (r *DIVReg) RetrieveTIMAIncrementRequest() bool {
	value := r.timaIncReq
	r.timaIncReq = false
	return value
}

func (r *DIVReg) handleIncrementRequest() {
	bitValue := r.tac.Enabled.Get()
	bitValue &= uint8(bitwise.GetBit16(r.cyclesCounter, _divBitForTac[r.tac.Rate.Get()]))

	if r.previousBitValue == 1 && bitValue == 0 {
		r.timaIncReq = true
	}
	r.previousBitValue = bitValue
}
