package bitwise

func GetBits8(value, offset, mask uint8) uint8 {
	return (value >> offset) & mask
}

func GetBit8(value, offset uint8) uint8 {
	return GetBits8(value, offset, 0x1)
}

func GetBits16(value, offset, mask uint16) uint16 {
	return (value >> offset) & mask
}

func GetBit16(value, offset uint16) uint16 {
	return GetBits16(value, offset, 0x1)
}
