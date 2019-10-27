package misc

import opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"

func (s *unitTestSuite) TestHalt() {
	s.Regs.HaltMode = false

	result := s.factory.Halt()()
	s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
	s.True(s.Regs.HaltMode)
}
