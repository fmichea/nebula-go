package bitwise

func InverseMask8(value, mask uint8) uint8 {
	return (value | mask) ^ mask
}

func InverseMask16(value, mask uint16) uint16 {
	return (value | mask) ^ mask
}

func InverseMask32(value, mask uint32) uint32 {
	return (value | mask) ^ mask
}
