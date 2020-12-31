package misc

import (
	"fmt"

	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
)

func (s *unitTestSuite) TestSCF() {
	cases := []struct {
		initialFlags uint8
		resultFlags  uint8
	}{
		{registers.FlagsCleared, registers.CY},
		{registers.FlagsFullSet, registers.ZF | registers.CY},
	}

	fn := s.factory.SCF()

	for _, c := range cases {
		name := fmt.Sprintf("scf with flags (initial = %#v, result = %#v)", c.initialFlags, c.resultFlags)

		s.Run(name, func() {
			s.Regs.F.Set(c.initialFlags)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
			s.EqualFlags(c.resultFlags)
		})
	}
}
