package cartridge

import (
	"nebula-go/pkg/gbc/memory/lib"
)

const (
	_cgbFlagAddress             = 0x143
	_cgbFlagBackwardsCompatible = 0x80
	_cgbFlagColorOnly           = 0xC0
)

func loadROMType(romData []uint8) lib.ROMType {
	value := romData[_cgbFlagAddress]
	if value == _cgbFlagBackwardsCompatible || value == _cgbFlagColorOnly {
		return lib.CGB001
	}
	return lib.DMG01
}
