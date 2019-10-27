package registers

type Byte interface {
	Get() uint8
	Set(value uint8)
}

// TODO: Maybe refactor this with the Options pattern? Although this might collide with
//  doing options for the DByte, BitProxy and Flag also.

func NewByte(value uint8) Byte {
	return NewByteWithMask(value, 0xFF)
}

func NewByteWithMask(value, mask uint8) Byte {
	return NewMappedByteWithMask(&value, mask)
}

func NewMappedByte(valuePtr *uint8) Byte {
	return NewMappedByteWithMask(valuePtr, 0xFF)
}

func NewMappedByteWithMask(valuePtr *uint8, mask uint8) Byte {
	*valuePtr &= mask
	return &byte{
		valuePtr: valuePtr,
		mask:     mask,
	}
}

type byte struct {
	valuePtr *uint8
	mask     uint8
}

func (b *byte) Set(value uint8) {
	*b.valuePtr = value & b.mask
}

func (b *byte) Get() uint8 {
	return *b.valuePtr
}
