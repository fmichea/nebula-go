package cb

import (
	"fmt"
	"nebula-go/pkg/common/testhelpers"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) TestCB_ValidCase() {
	cases := []struct {
		aValue uint8

		initialFlags uint8
		resultFlags  uint8
	}{
		{0x82, z80lib.FlagsCleared, z80lib.ZF | z80lib.HC},
		{0x82, z80lib.FlagsFullSet, z80lib.ZF | z80lib.HC | z80lib.CY},
		{0x43, z80lib.FlagsCleared, z80lib.HC},
		{0x43, z80lib.FlagsFullSet, z80lib.HC | z80lib.CY},
	}

	fn := s.factory.CB()

	for _, c := range cases {
		name := fmt.Sprintf("cb opcode 0x47 (bit 0, %%a) with a = %#v", c.aValue)

		s.Run(name, func() {
			pc := s.Regs.PC

			s.Regs.A.Set(c.aValue)
			s.Regs.F.Set(c.initialFlags)

			s.MockMMU.EXPECT().ReadByte(pc+1).Return(uint8(0x47), nil)

			result := fn()
			s.Equal(opcodeslib.OpcodeSuccess(2, 8), result)
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestCB_InvalidRead() {
	pc := s.Regs.PC

	s.MockMMU.EXPECT().ReadByte(pc+1).Return(uint8(0x00), testhelpers.ErrTesting1)

	result := s.factory.CB()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
