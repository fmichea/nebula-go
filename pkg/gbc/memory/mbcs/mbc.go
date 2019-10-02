package mbcs

import (
	"nebula-go/pkg/gbc/memory/lib"
)

type MBC interface {
	ContainsAddress(addr uint16) bool
	BytePtr(accessType lib.AccessType, addr uint16, value uint8) (*uint8, error)
}
