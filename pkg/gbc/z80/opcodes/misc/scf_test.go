package misc

import (
	"fmt"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) TestSCF() {
	cases := []struct {
		initialFlags uint8
		resultFlags  uint8
	}{
		{z80lib.FlagsCleared, z80lib.CY},
		{z80lib.FlagsFullSet, z80lib.ZF | z80lib.CY},
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
