package cartridge

import (
	"nebula-go/pkg/gbc/memory/lib"
)

var (
	_targetMarketAddress = 0x14A
)

func loadROMMarket(romData []uint8) lib.ROMMarket {
	if romData[_targetMarketAddress] != 0 {
		return lib.NonJapanese
	}
	return lib.Japanese
}
