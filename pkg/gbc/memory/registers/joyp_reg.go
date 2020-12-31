package registers

import (
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

type JOYPButton struct {
	registerslib.Flag
}

func NewJOYPButton(reg registerslib.Byte, bit uint8) *JOYPButton {
	return &JOYPButton{
		Flag: registerslib.NewFlag(reg, bit),
	}
}

func (b *JOYPButton) Press() {
	b.SetBool(false)
}

func (b *JOYPButton) Release() {
	b.SetBool(true)
}

type JOYPReg struct {
	registerslib.Byte

	ifReg *InterruptReg

	selectorReg   registerslib.Byte
	directionsReg registerslib.Byte
	buttonsReg    registerslib.Byte

	ButtonKeys    registerslib.Flag
	DirectionKeys registerslib.Flag

	DownButton  *JOYPButton
	UpButton    *JOYPButton
	LeftButton  *JOYPButton
	RightButton *JOYPButton

	StartButton  *JOYPButton
	SelectButton *JOYPButton
	BButton      *JOYPButton
	AButton      *JOYPButton
}

func NewJOYPReg(ptr *uint8) *JOYPReg {
	var selectorRegValue uint8
	selectorReg := registerslib.NewByteWithMask(&selectorRegValue, 0xF0, 0x30)

	directionsReg := registerslib.NewThreadSafeByteWithMask(0x0F, 0x0F)
	buttonsReg := registerslib.NewThreadSafeByteWithMask(0x0F, 0x0F)

	return &JOYPReg{
		Byte: selectorReg,

		selectorReg:   selectorReg,
		directionsReg: directionsReg,
		buttonsReg:    buttonsReg,

		ButtonKeys:    registerslib.NewFlag(selectorReg, 5),
		DirectionKeys: registerslib.NewFlag(selectorReg, 4),

		DownButton:  NewJOYPButton(directionsReg, 3),
		UpButton:    NewJOYPButton(directionsReg, 2),
		LeftButton:  NewJOYPButton(directionsReg, 1),
		RightButton: NewJOYPButton(directionsReg, 0),

		StartButton:  NewJOYPButton(buttonsReg, 3),
		SelectButton: NewJOYPButton(buttonsReg, 2),
		BButton:      NewJOYPButton(buttonsReg, 1),
		AButton:      NewJOYPButton(buttonsReg, 0),
	}
}

func (r *JOYPReg) Get() uint8 {
	return r.selectorReg.Get() | r.getLower4()
}

func (r *JOYPReg) getLower4() uint8 {
	if !r.DirectionKeys.GetBool() {
		return r.directionsReg.Get()
	} else if !r.ButtonKeys.GetBool() {
		return r.buttonsReg.Get()
	}
	return 0x0F
}
