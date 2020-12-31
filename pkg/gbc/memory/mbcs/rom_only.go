package mbcs

import (
	"nebula-go/pkg/gbc/memory/lib"
	"nebula-go/pkg/gbc/memory/segments"
)

type romOnly struct {
	rom  segments.Segment
	eram segments.Segment
}

func newRomOnly(rom, eram segments.Segment) MBC {
	return &romOnly{
		rom:  rom,
		eram: eram,
	}
}

func (r *romOnly) AddressRanges() []segments.AddressRange {
	var result []segments.AddressRange

	result = append(result, r.rom.AddressRanges()...)
	result = append(result, r.eram.AddressRanges()...)
	return result
}

func (r *romOnly) ContainsAddress(addr uint16) bool {
	return r.rom.ContainsAddress(addr) || r.eram.ContainsAddress(addr)
}

func (r *romOnly) ReadByte(addr uint16) (uint8, error) {
	if r.rom.ContainsAddress(addr) {
		return r.rom.ReadByte(addr)
	}
	return 0, lib.ErrInvalidRead
}

func (r *romOnly) ReadByteSlice(addr uint16, count uint) ([]uint8, error) {
	if r.rom.ContainsAddress(addr) {
		return r.rom.ReadByteSlice(addr, count)
	}
	return nil, lib.ErrInvalidRead
}

func (r *romOnly) WriteByte(addr uint16, value uint8) error {
	return lib.ErrInvalidWrite
}

func (r *romOnly) WriteByteSlice(addr uint16, values []uint8) error {
	// TODO: there is nothing in specification disallowing this, not currently used but could be implemented.
	return ErrMBCSliceOperatioInvalid
}

func (r *romOnly) ByteHook(addr uint16) (*uint8, error) {
	return nil, ErrMBCHookInvalid
}
