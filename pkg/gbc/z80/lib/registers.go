package z80lib

import (
	"nebula-go/pkg/gbc/memory/registers"
)

type Registers struct {
	A registers.Byte
	F *Flags
	B registers.Byte
	C registers.Byte
	D registers.Byte
	E registers.Byte
	H registers.Byte
	L registers.Byte

	AF registers.DByte
	BC registers.DByte
	DE registers.DByte
	HL registers.DByte

	SP registers.DByte
	PC uint16

	IME      bool
	HaltMode bool
}

func NewRegisters() *Registers {
	a := registers.NewByte(0x11)
	f := NewFlags(0xB0)

	b := registers.NewByte(0x00)
	c := registers.NewByte(0x13)

	d := registers.NewByte(0x00)
	e := registers.NewByte(0xD8)

	h := registers.NewByte(0x01)
	l := registers.NewByte(0x4D)

	return &Registers{
		A: a,
		F: f,
		B: b,
		C: c,
		D: d,
		E: e,
		H: h,
		L: l,

		AF: registers.NewSplitDByte(a, f),
		BC: registers.NewSplitDByte(b, c),
		DE: registers.NewSplitDByte(d, e),
		HL: registers.NewSplitDByte(h, l),

		PC: 0x0100,
		SP: registers.NewDByte(0xFFFE),

		// FIXME: initialize IME flag.
	}
}
