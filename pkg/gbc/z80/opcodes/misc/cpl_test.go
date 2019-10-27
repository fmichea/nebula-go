package misc

import (
	"fmt"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) TestCPL() {
	cases := []struct {
		initialValue uint8
		resultValue  uint8

		initialFlags uint8
		resultFlags  uint8
	}{
		{0x00, 0xFF, z80lib.FlagsCleared, z80lib.HC | z80lib.NE},
		{0x00, 0xFF, z80lib.FlagsFullSet, z80lib.FlagsFullSet},

		{0xF0, 0x0F, z80lib.FlagsCleared, z80lib.HC | z80lib.NE},
		{0xF0, 0x0F, z80lib.FlagsFullSet, z80lib.FlagsFullSet},

		{0x3C, 0xC3, z80lib.FlagsCleared, z80lib.HC | z80lib.NE},
		{0xC3, 0x3C, z80lib.FlagsFullSet, z80lib.FlagsFullSet},
	}

	fn := s.factory.CPL()

	for _, c := range cases {
		name := fmt.Sprintf(
			"cpl test a (initial = %#v, result = %#v) with flags (initial = %#v, result = %#v)",
			c.initialValue,
			c.resultValue,
			c.initialFlags,
			c.resultFlags,
		)

		s.Run(name, func() {
			s.Regs.A.Set(c.initialValue)
			s.Regs.F.Set(c.initialFlags)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
			s.Equal(c.resultValue, s.Regs.A.Get())
			s.EqualFlags(c.resultFlags)
		})
	}
}
