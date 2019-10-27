package load

import (
	"nebula-go/pkg/common/testhelpers"
	"nebula-go/pkg/gbc/memory/registers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) TestHLPtrToByte_ValidCase() {
	reg := registers.NewByte(0x00)

	value := uint8(0xDE)
	hl := uint16(0xABCD)

	s.Regs.HL.Set(hl)

	s.MockMMU.EXPECT().ReadByte(hl).Return(value, nil)

	result := s.factory.HLPtrToByte(reg)()
	s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
	s.Equal(value, reg.Get())
}

func (s *unitTestSuite) TestHLPtrToByte_InvalidRead() {
	reg := registers.NewByte(0x00)

	value := uint8(0xDE)
	hl := uint16(0xABCD)

	s.Regs.HL.Set(hl)

	s.MockMMU.EXPECT().ReadByte(hl).Return(value, testhelpers.ErrTesting1)

	result := s.factory.HLPtrToByte(reg)()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestByteToHLPtr_ValidCase() {
	value := uint8(0xDE)
	hl := uint16(0xABCD)

	s.Regs.HL.Set(hl)

	reg := registers.NewByte(value)

	s.MockMMU.EXPECT().WriteByte(hl, value).Return(nil)

	result := s.factory.ByteToHLPtr(reg)()
	s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
}

func (s *unitTestSuite) TestByteToHLPtr_InvalidWrite() {
	value := uint8(0xDE)
	hl := uint16(0xABCD)

	s.Regs.HL.Set(hl)

	reg := registers.NewByte(value)

	s.MockMMU.EXPECT().WriteByte(hl, value).Return(testhelpers.ErrTesting1)

	result := s.factory.ByteToHLPtr(reg)()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestD8ToHLPtr_ValidCase() {
	value := uint8(0xDE)
	hl := uint16(0xABCD)

	s.Regs.HL.Set(hl)

	s.MockMMU.EXPECT().ReadByte(s.Regs.PC+1).Return(value, nil)
	s.MockMMU.EXPECT().WriteByte(hl, value).Return(nil)

	result := s.factory.D8ToHLPtr()()
	s.Equal(opcodeslib.OpcodeSuccess(2, 12), result)
}

func (s *unitTestSuite) TestD8ToHLPtr_InvalidRead() {
	value := uint8(0xDE)
	hl := uint16(0xABCD)

	s.Regs.HL.Set(hl)

	s.MockMMU.EXPECT().ReadByte(s.Regs.PC+1).Return(value, testhelpers.ErrTesting1)

	result := s.factory.D8ToHLPtr()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestD8ToHLPtr_InvalidWrite() {
	value := uint8(0xDE)
	hl := uint16(0xABCD)

	s.Regs.HL.Set(hl)

	s.MockMMU.EXPECT().ReadByte(s.Regs.PC+1).Return(value, nil)
	s.MockMMU.EXPECT().WriteByte(hl, value).Return(testhelpers.ErrTesting1)

	result := s.factory.D8ToHLPtr()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
