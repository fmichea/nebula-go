package bitwise

func ConvertHighLow8To16(high, low uint8) uint16 {
	return uint16(high)<<8 | uint16(low)
}
