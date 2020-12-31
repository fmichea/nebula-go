package registers

import (
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

type KEY1CurrentSpeedFlag struct {
	registerslib.Flag

	value bool
}

func NewKEY1CurrentSpeedFlag(reg registerslib.Byte) *KEY1CurrentSpeedFlag {
	return &KEY1CurrentSpeedFlag{
		Flag: registerslib.NewFlag(reg, 7),
	}
}

func (f *KEY1CurrentSpeedFlag) SetBool(value bool) {
	f.Flag.SetBool(value)
	f.value = value
}

func (f *KEY1CurrentSpeedFlag) IsDoubleSpeed() bool {
	return f.value
}

type KEY1Reg struct {
	registerslib.Byte

	CurrentSpeed  *KEY1CurrentSpeedFlag
	ChangeRequest registerslib.Flag
}

func NewKEY1Reg(ptr *uint8) *KEY1Reg {
	reg := registerslib.NewByteWithMask(ptr, 0x00, 0x01)

	return &KEY1Reg{
		Byte: reg,

		CurrentSpeed:  NewKEY1CurrentSpeedFlag(reg),
		ChangeRequest: registerslib.NewFlag(reg, 0),
	}
}

func (r *KEY1Reg) SwitchIfRequested() {
	if r.ChangeRequest.GetBool() {
		r.CurrentSpeed.SetBool(!r.CurrentSpeed.GetBool())
		r.ChangeRequest.SetBool(false)
	}
}
