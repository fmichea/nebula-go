package cartridge

var (
	_checksumStartAddress = 0x134
	_checksumEndAddress   = 0x14D
)

func verifyHeaderChecksum(romData []uint8) (err error) {
	checksum := uint8(0)

	for addr := _checksumStartAddress; addr <= _checksumEndAddress; addr++ {
		checksum -= romData[addr] + 1
	}

	if checksum != 0xFF {
		return ErrChecksumInvalid
	}
	return nil
}
