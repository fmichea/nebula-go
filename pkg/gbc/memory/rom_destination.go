package memory

type ROMDestination int

const (
	Japanese ROMDestination = iota
	NonJapanese
)

func (t ROMDestination) String() string {
	switch t {
	case Japanese:
		return "Japanese"

	case NonJapanese:
		return "Non-Japanese"

	default:
		return "Unknown"
	}
}
