package memory

type ROMSize uint8

const (
	ROMSize32KB  ROMSize = 0x00 // no ROM banking
	ROMSize64KB  ROMSize = 0x01 // 4 banks
	ROMSize128KB ROMSize = 0x02 // 8 banks
	ROMSize256KB ROMSize = 0x03 // 16 banks
	ROMSize512KB ROMSize = 0x04 // 32 banks
	ROMSize1MB   ROMSize = 0x05 // 64 banks - only 63 banks used by MBC1
	ROMSize2MB   ROMSize = 0x06 // 128 banks - only 125 banks used by MBC1
	ROMSize4MB   ROMSize = 0x07 // 256 banks
	ROMSize1p1MB ROMSize = 0x52 // 72 banks
	ROMSize1p2MB ROMSize = 0x53 // 80 banks
	ROMSize1p5MB ROMSize = 0x54 // 96 banks

)

func (s ROMSize) String() string {
	values := map[ROMSize]string{
		ROMSize32KB:  "32KByte",
		ROMSize64KB:  "64KByte",
		ROMSize128KB: "128KByte",
		ROMSize256KB: "256KByte",
		ROMSize512KB: "512KByte",
		ROMSize1MB:   "1MByte",
		ROMSize2MB:   "2MByte",
		ROMSize4MB:   "4MByte",
		ROMSize1p1MB: "1.1MByte",
		ROMSize1p2MB: "1.2MByte",
		ROMSize1p5MB: "1.5MByte",
	}

	if name, ok := values[s]; ok {
		return name
	} else {
		return ""
	}
}

func (s ROMSize) IsValid() bool {
	return s.String() != ""
}
