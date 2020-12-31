package misc

import (
	"fmt"

	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
)

func (s *unitTestSuite) TestCPL() {
	cases := []struct {
		initialValue uint8
		resultValue  uint8

		initialFlags uint8
		resultFlags  uint8
	}{
		{0x00, 0xFF, registers.FlagsCleared, registers.HC | registers.NE},
		{0x00, 0xFF, registers.FlagsFullSet, registers.FlagsFullSet},

		{0xF0, 0x0F, registers.FlagsCleared, registers.HC | registers.NE},
		{0xF0, 0x0F, registers.FlagsFullSet, registers.FlagsFullSet},

		{0x3C, 0xC3, registers.FlagsCleared, registers.HC | registers.NE},
		{0xC3, 0x3C, registers.FlagsFullSet, registers.FlagsFullSet},
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
