package bitwise

func SetBits8(current, offset, mask, value uint8) uint8 {
	value = (value & mask) << offset

	current &= 0xFF ^ (mask << offset)
	current |= value
	return current
}
