package memory

func (m *mmu) WriteByte(addr uint16, value uint8) error {
	if err := m.writeByteInternal(addr, value); err != nil {
		return wrapWriteError(addr, err)
	}
	return nil
}

func (m *mmu) WriteDByte(addr, value uint16) error {
	if err := m.writeByteInternal(addr, uint8(value&0xFF)); err != nil {
		return wrapWriteError(addr, err)
	}
	if err := m.writeByteInternal(addr+1, uint8(value>>8)); err != nil {
		return wrapWriteError(addr, err)
	}
	return nil
}

func (m *mmu) WriteByteSlice(addr uint16, values []uint8) error {
	segment, err := m.getSegmentForAddress(addr)
	if err != nil {
		return wrapWriteError(addr, err)
	}

	if err := segment.WriteByteSlice(addr, values); err != nil {
		return wrapWriteError(addr, err)
	}
	return nil
}

func (m *mmu) writeByteInternal(addr uint16, value uint8) error {
	if hook, ok := m.hooks[addr]; ok {
		return hook.Set(value)
	}

	segment, err := m.getSegmentForAddress(addr)
	if err != nil {
		return err
	}

	return segment.WriteByte(addr, value)
}
