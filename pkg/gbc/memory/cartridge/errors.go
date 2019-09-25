package cartridge

import (
	"errors"
)

var (
	ErrInvalidROMRead      = errors.New("failed to read ROM in full")
	ErrChecksumInvalid     = errors.New("cartridge checksum is invalid")
	ErrNintendoLogoInvalid = errors.New("nintendo logo from ROM did not match")
	ErrROMSizeInvalid      = errors.New("got an invalid ROM size value")
	ErrRAMSizeInvalid      = errors.New("got an invalid RAM size value")
)
