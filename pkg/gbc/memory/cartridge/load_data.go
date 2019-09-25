package cartridge

import (
	"os"
)

func loadData(filename string) ([]uint8, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	size := int(stat.Size())

	romBuffer := make([]uint8, size)

	readCount, err := f.Read(romBuffer)
	if err != nil {
		return nil, err
	}

	if readCount != size {
		return nil, ErrInvalidROMRead
	}

	return romBuffer, nil
}
