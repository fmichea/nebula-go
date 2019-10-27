package load

import (
	"nebula-go/pkg/common/testhelpers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) TestAToHighRAM_ValidCase() {
	a := uint8(0xDE)

	offset := uint8(0xBC)
	pc := s.Regs.PC

	s.Regs.A.Set(a)

	s.MockMMU.EXPECT().ReadByte(pc+1).Return(offset, nil)
	s.MockMMU.EXPECT().WriteByte(0xFF00+uint16(offset), a).Return(nil)

	result := s.factory.AToHighRAM()()
	s.Equal(opcodeslib.OpcodeSuccess(2, 12), result)
}

func (s *unitTestSuite) TestAToHighRAM_InvalidRead() {
	a := uint8(0xDE)

	offset := uint8(0xBC)
	pc := s.Regs.PC

	s.Regs.A.Set(a)

	s.MockMMU.EXPECT().ReadByte(pc+1).Return(offset, testhelpers.ErrTesting1)

	result := s.factory.AToHighRAM()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestAToHighRAM_InvalidWrite() {
	a := uint8(0xDE)

	offset := uint8(0xBC)
	pc := s.Regs.PC

	s.Regs.A.Set(a)

	s.MockMMU.EXPECT().ReadByte(pc+1).Return(offset, nil)
	s.MockMMU.EXPECT().WriteByte(0xFF00+uint16(offset), a).Return(testhelpers.ErrTesting1)

	result := s.factory.AToHighRAM()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestHighRAMToA_ValidCase() {
	value := uint8(0xDE)

	offset := uint8(0xBC)
	pc := s.Regs.PC

	s.Regs.A.Set(0x00)

	s.MockMMU.EXPECT().ReadByte(pc+1).Return(offset, nil)
	s.MockMMU.EXPECT().ReadByte(0xFF00+uint16(offset)).Return(value, nil)

	result := s.factory.HighRAMToA()()
	s.Equal(opcodeslib.OpcodeSuccess(2, 12), result)
	s.Equal(value, s.Regs.A.Get())
}

func (s *unitTestSuite) TestHighRAMToA_InvalidFirstRead() {
	offset := uint8(0xBC)
	pc := s.Regs.PC

	s.Regs.A.Set(0x00)

	s.MockMMU.EXPECT().ReadByte(pc+1).Return(offset, testhelpers.ErrTesting1)

	result := s.factory.HighRAMToA()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestHighRAMToA_InvalidSecondRead() {
	value := uint8(0xDE)

	offset := uint8(0xBC)
	pc := s.Regs.PC

	s.Regs.A.Set(0x00)

	s.MockMMU.EXPECT().ReadByte(pc+1).Return(offset, nil)
	s.MockMMU.EXPECT().ReadByte(0xFF00+uint16(offset)).Return(value, testhelpers.ErrTesting1)

	result := s.factory.HighRAMToA()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestAToCPtrInHighRAM_ValidCase() {
	a := uint8(0xDE)
	offset := uint8(0xBC)

	s.Regs.C.Set(offset)
	s.Regs.A.Set(a)

	s.MockMMU.EXPECT().WriteByte(0xFF00+uint16(offset), a).Return(nil)

	result := s.factory.AToCPtrInHighRAM()()
	s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
}

func (s *unitTestSuite) TestAToCPtrInHighRAM_InvalidWrite() {
	a := uint8(0xDE)
	offset := uint8(0xBC)

	s.Regs.C.Set(offset)
	s.Regs.A.Set(a)

	s.MockMMU.EXPECT().WriteByte(0xFF00+uint16(offset), a).Return(testhelpers.ErrTesting1)

	result := s.factory.AToCPtrInHighRAM()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestCPtrInHighRAMToA_ValidCase() {
	value := uint8(0xDE)
	offset := uint8(0xBC)

	s.Regs.C.Set(offset)
	s.Regs.A.Set(0x00)

	s.MockMMU.EXPECT().ReadByte(0xFF00+uint16(offset)).Return(value, nil)

	result := s.factory.CPtrInHighRAMToA()()
	s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
	s.Equal(value, s.Regs.A.Get())
}

func (s *unitTestSuite) TestCPtrInHighRAMToA_InvalidSecondRead() {
	value := uint8(0xDE)
	offset := uint8(0xBC)

	s.Regs.C.Set(offset)
	s.Regs.A.Set(0x00)

	s.MockMMU.EXPECT().ReadByte(0xFF00+uint16(offset)).Return(value, testhelpers.ErrTesting1)

	result := s.factory.CPtrInHighRAMToA()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
