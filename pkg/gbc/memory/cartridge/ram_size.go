package cartridge

import (
	"nebula-go/pkg/gbc/memory/lib"
)

var (
	_ramSizeAddress = 0x149
)

func loadRAMSize(romData []uint8) lib.RAMSize {
	return lib.RAMSize(romData[_ramSizeAddress])
}

func verifyRAMSizeFlag(romData []uint8) error {
	if !loadRAMSize(romData).IsValid() {
		return ErrRAMSizeInvalid
	}
	return nil
}
