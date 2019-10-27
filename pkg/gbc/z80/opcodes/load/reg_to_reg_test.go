package load

import (
	"nebula-go/pkg/gbc/memory/registers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) TestByteToByte() {
	v8 := uint8(0xAB)

	reg1 := registers.NewByte(0x00)
	reg2 := registers.NewByte(v8)

	result := s.factory.ByteToByte(reg1, reg2)()

	s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
	s.Equal(v8, reg1.Get())
	s.Equal(v8, reg2.Get())
}
