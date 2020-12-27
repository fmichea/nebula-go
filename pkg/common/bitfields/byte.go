package bitfields

type ROByte interface {
	Get() uint8
}

type Byte interface {
	Get() uint8
	Set(value uint8)
}

func NewByte(value uint8) *byteReg {
	return NewByteWithMask(value, 0xFF)
}

func NewByteWithMask(value, mask uint8) *byteReg {
	return &byteReg{
		value: value & mask,
		mask:  mask,
	}
}

type byteReg struct {
	value uint8
	mask  uint8
}

func (b *byteReg) Set(value uint8) {
	b.value = value & b.mask
}

func (b *byteReg) Get() uint8 {
	return b.value
}
