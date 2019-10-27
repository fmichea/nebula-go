package misc

import opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"

func (s *unitTestSuite) TestStop() {
	result := s.factory.Stop()()
	s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
}
