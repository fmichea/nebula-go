package misc

import (
	"fmt"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) TestRRCA() {
	cases := []struct {
		initialValue uint8
		resultValue  uint8

		initialFlags uint8
		resultFlags  uint8
	}{
		{0x00, 0x00, z80lib.FlagsCleared, z80lib.FlagsCleared},
		{0x00, 0x00, z80lib.FlagsFullSet, z80lib.ZF},
		{0xF0, 0x78, z80lib.FlagsCleared, z80lib.FlagsCleared},
		{0xF1, 0xF8, z80lib.FlagsCleared, z80lib.CY},
		{0xF1, 0xF8, z80lib.FlagsFullSet, z80lib.ZF | z80lib.CY},
		{0x01, 0x80, z80lib.FlagsCleared, z80lib.CY},
	}

	fn := s.factory.RRCA()

	for _, c := range cases {
		name := fmt.Sprintf(
			"rrca test a (initial = %#v, result = %#v) with flags (initial = %#v, result = %#v)",
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
