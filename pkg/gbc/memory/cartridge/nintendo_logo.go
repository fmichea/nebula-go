package cartridge

import (
	"bytes"
)

var (
	NintendoLogo = []uint8(
		"\xce\xed\x66\x66\xcc\x0d\x00\x0b\x03\x73\x00\x83\x00\x0c\x00\x0d" +
			"\x00\x08\x11\x1f\x88\x89\x00\x0e\xdc\xcc\x6e\xe6\xdd\xdd\xd9\x99" +
			"\xbb\xbb\x67\x63\x6e\x0e\xec\xcc\xdd\xdc\x99\x9f\xbb\xb9\x33\x3e")

	_nintendoLogoStartAddress = 0x104
	_nintendoLogoEndAddress   = _nintendoLogoStartAddress + len(NintendoLogo)
)

func verifyNintendoLogo(romData []uint8) error {
	nintendoLogo := romData[_nintendoLogoStartAddress:_nintendoLogoEndAddress]
	if !bytes.Equal(NintendoLogo, nintendoLogo) {
		return ErrNintendoLogoInvalid
	}
	return nil
}
