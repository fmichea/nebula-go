package load

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (s *unitTestSuite) TestByteToByte() {
	v8 := uint8(0xAB)

	reg1 := registerslib.NewByte(0x00)
	reg2 := registerslib.NewByte(v8)

	result := s.factory.ByteToByte(reg1, reg2)()

	s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
	s.Equal(v8, reg1.Get())
	s.Equal(v8, reg2.Get())
}
