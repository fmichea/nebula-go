package registers

type dbyteSplitByte struct {
	msb Byte
	lsb Byte
}

func NewSplitDByte(msb, lsb Byte) DByte {
	return &dbyteSplitByte{
		msb: msb,
		lsb: lsb,
	}
}

func (d *dbyteSplitByte) Set(value uint16) {
	d.msb.Set(uint8(value >> 8))
	d.lsb.Set(uint8(value & 0xFF))
}

func (d *dbyteSplitByte) Get() uint16 {
	var result uint16

	result |= uint16(d.msb.Get()) << 8
	result |= uint16(d.lsb.Get())
	return result
}
