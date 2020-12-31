package registers

import registerslib "nebula-go/pkg/gbc/memory/registers/lib"

type InterruptFlag struct {
	registerslib.Flag
}

func NewInterruptFlag(reg registerslib.Byte, bit uint8) *InterruptFlag {
	return &InterruptFlag{
		Flag: registerslib.NewFlag(reg, bit),
	}
}

func (f *InterruptFlag) IsRequested() bool {
	return f.GetBool()
}

func (f *InterruptFlag) Request() {
	f.SetBool(true)
}

func (f *InterruptFlag) Acknowledge() {
	f.SetBool(false)
}
