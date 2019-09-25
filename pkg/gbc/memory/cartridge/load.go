package cartridge

import (
	"io"
)

func Load(out io.Writer, filename string) (*ROM, error) {
	romData, err := loadData(filename)
	if err != nil {
		return nil, err
	}

	if err := validate(out, romData); err != nil {
		return nil, err
	}

	romType := loadROMType(romData)

	result := &ROM{
		Title: loadROMTitle(romData, romType),

		Type:   romType,
		Size:   loadROMSize(romData),
		Market: loadROMMarket(romData),

		RAMSize: loadRAMSize(romData),

		MBCSelector: loadMBCSelector(romData),

		Data: romData,
	}

	if err := result.PrintInformation(out); err != nil {
		return nil, err
	}

	return result, nil
}
