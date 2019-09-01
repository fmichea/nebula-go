package memory

type RAMSize uint8

const (
	RAMSizeNone RAMSize = 0x00
	RAMSize2KB  RAMSize = 0x01
	RAMSize8KB  RAMSize = 0x02
	RAMSize32KB RAMSize = 0x03 // 4 banks of 8 KByte each
)

func (s RAMSize) String() string {
	values := map[RAMSize]string{
		RAMSizeNone: "None",
		RAMSize2KB:  "2KByte",
		RAMSize8KB:  "8KByte",
		RAMSize32KB: "32KByte",
	}

	if name, ok := values[s]; ok {
		return name
	} else {
		return ""
	}
}

func (s RAMSize) IsValid() bool {
	return s.String() != ""
}
