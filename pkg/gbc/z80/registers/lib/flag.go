package registerslib

import (
	"nebula-go/pkg/common/bitfields"
)

type Flag bitfields.Flag

func NewFlag(reg Byte, bit uint8) Flag {
	return bitfields.NewFlag(reg, bit)
}
