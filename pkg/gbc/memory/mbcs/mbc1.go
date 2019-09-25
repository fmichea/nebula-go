package mbcs

import (
	"nebula-go/pkg/gbc/memory/lib"
)

type MBC1 struct {
}

func (m *MBC1) ContainsAddress(addr uint16) bool {
	return false
}

func (m *MBC1) BytePtr(accessType lib.AccessType, addr uint16, value uint8) *uint8 {
	return nil
}
