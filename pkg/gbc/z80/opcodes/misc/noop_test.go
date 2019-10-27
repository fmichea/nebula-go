package misc

import opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"

func (s *unitTestSuite) TestNoop() {
	result := s.factory.Noop()()
	s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
}
