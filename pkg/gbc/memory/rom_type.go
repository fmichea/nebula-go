package memory

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
