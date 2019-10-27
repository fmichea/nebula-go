package bitwise

func HighBit8(value uint8) uint8 {
	return (value >> 7) & 0x1
}

func LowBit8(value uint8) uint8 {
	return value & 0x1
}
