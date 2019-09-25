package cartridge

import (
	"bytes"

	"nebula-go/pkg/gbc/memory/lib"
)

var (
	_titleAddress = 0x134
)

func loadROMTitle(romData []uint8, romType lib.ROMType) string {
	titleBytes := romData[_titleAddress : _titleAddress+romType.GetTitleSize()]

	zeroByteIndex := bytes.IndexByte(titleBytes, 0)
	if zeroByteIndex == -1 {
		zeroByteIndex = len(titleBytes)
	}
	return string(titleBytes[:zeroByteIndex])
}
