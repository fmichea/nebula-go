package registerslib

type Flag interface {
	BitProxy

	GetBool() bool
	SetBool(value bool)
}

func NewFlag(reg Byte, bit uint8) Flag {
	return &flag{
		BitProxy: NewBitProxy(reg, bit, 0x1),
	}
}

type flag struct {
	BitProxy
}

func (f *flag) SetBool(value bool) {
	if value {
		f.Set(0x1)
	} else {
		f.Set(0x0)
	}
}

func (f *flag) GetBool() bool {
	return f.Get() == 0x1
}
