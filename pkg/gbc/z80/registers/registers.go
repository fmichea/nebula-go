package registers

import (
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

type Registers struct {
	A registerslib.Byte
	F *Flags
	B registerslib.Byte
	C registerslib.Byte
	D registerslib.Byte
	E registerslib.Byte
	H registerslib.Byte
	L registerslib.Byte

	AF registerslib.DByte
	BC registerslib.DByte
	DE registerslib.DByte
	HL registerslib.DByte

	SP registerslib.DByte
	PC uint16

	IME      bool // Interrupt Master Enable
	HaltMode bool
}

func New() *Registers {
	a := registerslib.NewByte(0x11)
	f := NewFlags(0xB0)

	b := registerslib.NewByte(0x00)
	c := registerslib.NewByte(0x13)

	d := registerslib.NewByte(0x00)
	e := registerslib.NewByte(0xD8)

	h := registerslib.NewByte(0x01)
	l := registerslib.NewByte(0x4D)

	return &Registers{
		A: a,
		F: f,
		B: b,
		C: c,
		D: d,
		E: e,
		H: h,
		L: l,

		AF: registerslib.NewSplitDByte(a, f),
		BC: registerslib.NewSplitDByte(b, c),
		DE: registerslib.NewSplitDByte(d, e),
		HL: registerslib.NewSplitDByte(h, l),

		PC: 0x0100,
		SP: registerslib.NewDByte(0xFFFE),

		// FIXME: initialize IME flag.
	}
}
