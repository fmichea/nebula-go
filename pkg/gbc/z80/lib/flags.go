package z80lib

import (
	"strings"

	"nebula-go/pkg/gbc/memory/registers"
)

const (
	ZF uint8 = 1 << 7
	NE       = 1 << 6
	HC       = 1 << 5
	CY       = 1 << 4
)

const (
	FlagsFullSet uint8 = ZF | NE | HC | CY
	FlagsCleared       = 0
)

type Flags struct {
	registers.Byte

	ZF registers.Flag
	NE registers.Flag
	HC registers.Flag
	CY registers.Flag
}

func NewFlags(value uint8) *Flags {
	reg := registers.NewByteWithMask(value, 0xF0)

	return &Flags{
		Byte: reg,

		ZF: registers.NewFlag(reg, 7),
		NE: registers.NewFlag(reg, 6),
		HC: registers.NewFlag(reg, 5),
		CY: registers.NewFlag(reg, 4),
	}
}

func (f *Flags) String() string {
	var result []string

	for _, v := range []struct {
		name string
		flag registers.Flag
	}{
		{"ZF", f.ZF},
		{"NE", f.NE},
		{"HC", f.HC},
		{"CY", f.CY},
	} {
		if v.flag.GetBool() {
			result = append(result, v.name)
		}
	}
	return strings.Join(result, " | ")
}
