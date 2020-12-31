package registers

import (
	"strings"

	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
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
	registerslib.Byte

	ZF registerslib.Flag
	NE registerslib.Flag
	HC registerslib.Flag
	CY registerslib.Flag
}

func NewFlags(value uint8) *Flags {
	reg := registerslib.NewByteWithMask(value, 0xF0)

	return &Flags{
		Byte: reg,

		ZF: registerslib.NewFlag(reg, 7),
		NE: registerslib.NewFlag(reg, 6),
		HC: registerslib.NewFlag(reg, 5),
		CY: registerslib.NewFlag(reg, 4),
	}
}

func (f *Flags) String() string {
	var result []string

	for _, v := range []struct {
		name string
		flag registerslib.Flag
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
