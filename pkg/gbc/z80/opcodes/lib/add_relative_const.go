package opcodeslib

import (
	"nebula-go/pkg/common/bitwise"
)

func AddRelativeConstForMaskFunc(addr uint16, d8 uint8) func(uint16) uint16 {
	return func(mask uint16) uint16 {
		r8 := int8(d8)

		if r8 < 0 {
			return bitwise.Mask16(addr, mask) - bitwise.Mask16(uint16(-1*r8), mask)
		} else {
			return bitwise.Mask16(addr, mask) + bitwise.Mask16(uint16(r8), mask)
		}
	}
}

func AddRelativeConst(addr uint16, d8 uint8) uint16 {
	return AddRelativeConstForMaskFunc(addr, d8)(0xFFFF)
}
