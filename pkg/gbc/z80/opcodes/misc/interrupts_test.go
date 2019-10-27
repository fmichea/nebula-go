package misc

import (
	"fmt"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) TestDI() {
	cases := []struct {
		initialIME bool
	}{
		{true},
		{false},
	}

	fn := s.factory.DI()

	for _, c := range cases {
		name := fmt.Sprintf("di with original IME = %t", c.initialIME)

		s.Run(name, func() {
			s.Regs.IME = c.initialIME

			result := fn()
			s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
			s.False(s.Regs.IME)
		})
	}
}

func (s *unitTestSuite) TestEI() {
	cases := []struct {
		initialIME bool
	}{
		{true},
		{false},
	}

	fn := s.factory.EI()

	for _, c := range cases {
		name := fmt.Sprintf("ei with original IME = %t", c.initialIME)

		s.Run(name, func() {
			s.Regs.IME = c.initialIME

			result := fn()
			s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
			s.True(s.Regs.IME)
		})
	}
}
