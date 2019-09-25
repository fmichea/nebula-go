package cartridge

import (
	"nebula-go/pkg/gbc/memory/lib"
)

var (
	_romSizeAddress = 0x148
)

func loadROMSize(romData []uint8) lib.ROMSize {
	return lib.ROMSize(romData[_romSizeAddress])
}

func verifyROMSizeFlag(romData []uint8) error {
	if !loadROMSize(romData).IsValid() {
		return ErrROMSizeInvalid
	}
	return nil
}
