package lib

type MemoryIO interface {
	ReadByte(addr uint16) (uint8, error)
	WriteByte(addr uint16, value uint8) error

	ReadDByte(addr uint16) (uint16, error)
	WriteDByte(addr, value uint16) error

	ReadByteSlice(addr uint16, count uint) ([]uint8, error)
	WriteByteSlice(addr uint16, values []uint8) error

	ByteHook(addr uint16, fn func(ptr *uint8) (hook Hook, err error)) error
}
