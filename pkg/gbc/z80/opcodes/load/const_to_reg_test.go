package load

import (
	"nebula-go/pkg/common/testhelpers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (s *unitTestSuite) TestConstToByte_ValidCase() {
	value := uint8(0xA2)

	reg := registerslib.NewByte(0x00)

	s.MockMMU.EXPECT().ReadByte(s.Regs.PC+1).Return(value, nil)

	fn := s.factory.ConstToByte(reg)
	result := fn()

	s.Equal(opcodeslib.OpcodeSuccess(2, 8), result)
	s.Equal(value, reg.Get())
}

func (s *unitTestSuite) TestConstToByte_InvalidRead() {
	reg := registerslib.NewByte(0x00)

	s.MockMMU.EXPECT().ReadByte(s.Regs.PC+1).Return(uint8(0x00), testhelpers.ErrTesting1)

	fn := s.factory.ConstToByte(reg)
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestConstToDByte_ValidCase() {
	value := uint16(0xABCD)

	reg := registerslib.NewDByte(0x0000)

	s.MockMMU.EXPECT().ReadDByte(s.Regs.PC+1).Return(value, nil)

	fn := s.factory.ConstToDByte(reg)
	result := fn()

	s.Equal(opcodeslib.OpcodeSuccess(3, 12), result)
	s.Equal(value, reg.Get())
}

func (s *unitTestSuite) TestConstToDByte_InvalidRead() {
	reg := registerslib.NewDByte(0x0000)

	s.MockMMU.EXPECT().ReadDByte(s.Regs.PC+1).Return(uint16(0x0000), testhelpers.ErrTesting1)

	fn := s.factory.ConstToDByte(reg)
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
