package registerslib

import (
	"nebula-go/pkg/common/bitwise"
)

type Byte interface {
	Get() uint8
	Set(value uint8)
	SetNoMask(value uint8)
}

func NewByte(ptr *uint8, value uint8) Byte {
	return NewByteWithMask(ptr, value, 0xFF)
}

func NewByteWithMask(ptr *uint8, value, mask uint8) Byte {
	reg := &byteReg{
		ptr:  ptr,
		mask: mask,
	}
	reg.SetNoMask(value)
	return reg
}

type byteReg struct {
	ptr  *uint8
	mask uint8
}

func (b *byteReg) Set(value uint8) {
	*b.ptr = setMasked(value, *b.ptr, b.mask)
}

func (b *byteReg) SetNoMask(value uint8) {
	*b.ptr = value
}

func (b *byteReg) Get() uint8 {
	return *b.ptr
}

func setMasked(value, currentValue, mask uint8) uint8 {
	return bitwise.Mask8(value, mask) | bitwise.InverseMask8(currentValue, mask)
}
