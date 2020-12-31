package misc

import (
	"fmt"

	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
)

func (s *unitTestSuite) TestRRCA() {
	cases := []struct {
		initialValue uint8
		resultValue  uint8

		initialFlags uint8
		resultFlags  uint8
	}{
		{0x00, 0x00, registers.FlagsCleared, registers.FlagsCleared},
		{0x00, 0x00, registers.FlagsFullSet, registers.FlagsCleared},
		{0xF0, 0x78, registers.FlagsCleared, registers.FlagsCleared},
		{0xF1, 0xF8, registers.FlagsCleared, registers.CY},
		{0xF1, 0xF8, registers.FlagsFullSet, registers.CY},
		{0x01, 0x80, registers.FlagsCleared, registers.CY},
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
