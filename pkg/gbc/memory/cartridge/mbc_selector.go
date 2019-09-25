package cartridge

import (
	"nebula-go/pkg/gbc/memory/lib"
	"nebula-go/pkg/gbc/memory/mbcs"
)

var (
	_mbcFlagAddress = 0x147
)

func loadMBCSelector(romData []uint8) *mbcs.Selector {
	return mbcs.NewSelector(romData[_mbcFlagAddress])
}

func verifyMBCSelector(romData []uint8) error {
	if !loadMBCSelector(romData).IsValid() {
		return lib.ErrMBCNotImplemented
	}
	return nil
}
