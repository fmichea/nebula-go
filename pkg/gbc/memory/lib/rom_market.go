package lib

type ROMMarket int

const (
	Japanese ROMMarket = iota
	NonJapanese
)

func (t ROMMarket) String() string {
	switch t {
	case Japanese:
		return "Japanese"

	case NonJapanese:
		return "Non-Japanese"

	default:
		return "Unknown"
	}
}
