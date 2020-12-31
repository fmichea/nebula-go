package opcodeslib

func AddRelativeConst(v16 uint16, d8 uint8) uint16 {
	r8 := int8(d8)

	if r8 < 0 {
		return v16 - uint16(-1*r8)
	} else {
		return v16 + uint16(r8)
	}
}
