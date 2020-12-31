package registerslib

import (
	"nebula-go/pkg/common/bitfields"
)

type Byte bitfields.Byte

func NewByte(value uint8) Byte {
	return bitfields.NewByte(value)
}

func NewByteWithMask(value, mask uint8) Byte {
	return bitfields.NewByteWithMask(value, mask)
}
