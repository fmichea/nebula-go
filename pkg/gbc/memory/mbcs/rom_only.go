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

func (r *romOnly) ContainsAddress(addr uint16) bool {
	return r.rom.ContainsAddress(addr) || r.eram.ContainsAddress(addr)
}

func (r *romOnly) BytePtr(accessType lib.AccessType, addr uint16, value uint8) (ptr *uint8, err error) {
	switch accessType {
	case lib.AccessTypeRead:
		if r.rom.ContainsAddress(addr) {
			ptr = r.rom.BytePtr(addr)
		} else {
			err = ErrInvalidRead
		}

	case lib.AccessTypeWrite:
		err = ErrInvalidWrite
	}

	return
}
