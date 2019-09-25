package mbcs

import (
	"nebula-go/pkg/gbc/memory/lib"
	"nebula-go/pkg/gbc/memory/segments"
)

type RomOnlyMBC struct {
	ROM segments.Segment
}

func (r *RomOnlyMBC) ContainsAddress(addr uint16) bool {
	return false
}

func (r *RomOnlyMBC) BytePtr(accessType lib.AccessType, addr uint16, value uint8) *uint8 {
	switch accessType {
	case lib.AccessTypeRead:
		if r.ROM.ContainsAddress(addr) {
			return r.ROM.BytePtr(addr)
		}
	}
	return nil
}
