package lib

type ROMType int

const (
	DMG01 ROMType = iota
	CGB001
)

func (t ROMType) String() string {
	switch t {
	case DMG01:
		return "GAME BOY (DMG-01)"

	case CGB001:
		return "GAME BOY COLOR (CGB-001)"

	default:
		return "UNKNOWN"
	}
}

func (t ROMType) VRAMBankCount() uint {
	if t == CGB001 {
		return 2
	}
	return 1
}

func (t ROMType) WRAMBankCount() uint {
	if t == CGB001 {
		return 8
	}
	return 2
}

func (t ROMType) GetTitleSize() int {
	if t == DMG01 {
		return 0x16
	}
	return 0x15
}
