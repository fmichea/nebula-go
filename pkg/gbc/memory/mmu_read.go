package memory

import "nebula-go/pkg/common/bitwise"

func (m *mmu) ReadByte(addr uint16) (uint8, error) {
	value, err := m.readByteInternal(addr)
	if err != nil {
		return 0, wrapReadError(addr, err)
	}
	return value, nil
}

func (m *mmu) ReadDByte(addr uint16) (uint16, error) {
	high, err := m.readByteInternal(addr + 1)
	if err != nil {
		return 0, wrapReadError(addr, err)
	}

	low, err := m.readByteInternal(addr)
	if err != nil {
		return 0, wrapReadError(addr, err)
	}

	return bitwise.ConvertHighLow8To16(high, low), nil
}

func (m *mmu) ReadByteSlice(addr uint16, count uint) ([]uint8, error) {
	segment, err := m.getSegmentForAddress(addr)
	if err != nil {
		return nil, wrapReadError(addr, err)
	}

	values, err := segment.ReadByteSlice(addr, count)
	if err != nil {
		return nil, wrapReadError(addr, err)
	}
	return values, nil
}

func (m *mmu) readByteInternal(addr uint16) (uint8, error) {
	if hook, ok := m.hooks[addr]; ok {
		return hook.Get()
	}

	segment, err := m.getSegmentForAddress(addr)
	if err != nil {
		return 0, err
	}
	return segment.ReadByte(addr)
}
