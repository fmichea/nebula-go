package mbcs

import (
	"nebula-go/pkg/gbc/memory/segments"
)

var (
	_names = []string{
		/* 0x00 */
		"ROM ONLY",
		/* 0x01 -> 0x03: MBC1 */
		"MBC1",
		"MBC1+RAM",
		"MBC1+RAM+BATTERY",
		/* 0x04: Unused */
		"",
		/* 0x05 -> 0x06: MBC2 */
		"MBC2",
		"MBC2+BATTERY",
		/* 0x07: Unused */
		"",
		/* 0x08 -> 0x09: ROM+RAM */
		"ROM+RAM",
		"ROM+RAM+BATTERY",
		/* 0x0A: Unused */
		"",
		/* 0x0B -> 0x0D: MMM01 */
		"MMM01",
		"MMM01+RAM",
		"MMM01+RAM+BATTERY",
		/* 0x0E: Unused */
		"",
		/* 0x0F -> 0x13: MBC3 */
		"MBC3+TIMER+BATTERY",
		"MBC3+TIMER+RAM+BATTERY",
		"MBC3",
		"MBC3+RAM",
		"MBC3+RAM+BATTERY",
		/* 0x14: Unused */
		"",
		/* 0x15 -> 0x17: MBC4 */
		"MBC4",
		"MBC4+RAM",
		"MBC4+RAM+BATTERY",
		/* 0x18: Unused */
		"",
		/* 0x19 -> 0x1E: MBC5 */
		"MBC5",
		"MBC5+RAM",
		"MBC5+RAM+BATTERY",
		"MBC5+RUMBLE+RAM",
		"MBC5+RUMBLE+RAM+BATTERY",
	}

	_namesCount = uint8(len(_names))
)

type Selector struct {
	value uint8
}

func NewSelector(value uint8) *Selector {
	return &Selector{value: value}
}

func (s *Selector) IsValid() bool {
	return s.Name() != ""
}

func (s *Selector) GetMBC(rom segments.Segment, eram segments.Segment) MBC {
	switch s.value {
	case 0x00:
		return newRomOnly(rom, eram)

	case 0x01, 0x02, 0x03:
		return newMBC1(rom, eram)

	case 0x05, 0x06:
		return NewMBC2(rom, eram)

	case 0x0F, 0x10, 0x11, 0x12, 0x13:
	// load MBC3

	case 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E:
		return NewMBC5(rom, eram)

	default:
		// log bad thing?
	}
	return nil
}

func (s *Selector) Name() string {
	if s.value < _namesCount {
		return _names[s.value]
	} else {
		switch s.value {
		case 0xFC:
			return "POCKET CAMERA"
		case 0xFD:
			return "BANDAI TAMA5"
		case 0xFE:
			return "HuC3"
		case 0xFF:
			return "HuC1+RAM+BATTERY"
		}
	}
	return ""
}
