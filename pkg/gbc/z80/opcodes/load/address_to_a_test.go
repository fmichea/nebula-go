package load

import (
	"nebula-go/pkg/common/testhelpers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) TestAddressToA_ValidCase() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(0x00)

	s.MockMMU.EXPECT().ReadDByte(s.Regs.PC+1).Return(addr, nil)
	s.MockMMU.EXPECT().ReadByte(addr).Return(value, nil)

	result := s.factory.AddressToA()()
	s.Equal(opcodeslib.OpcodeSuccess(3, 16), result)
	s.Equal(value, s.Regs.A.Get())
}

func (s *unitTestSuite) TestAddressToA_InvalidFirstRead() {
	addr := uint16(0xABCD)

	s.Regs.A.Set(0x00)

	s.MockMMU.EXPECT().ReadDByte(s.Regs.PC+1).Return(addr, testhelpers.ErrTesting1)

	result := s.factory.AddressToA()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestAddressToA_InvalidSecondRead() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(0x00)

	s.MockMMU.EXPECT().ReadDByte(s.Regs.PC+1).Return(addr, nil)
	s.MockMMU.EXPECT().ReadByte(addr).Return(value, testhelpers.ErrTesting1)

	result := s.factory.AddressToA()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestBCPtrToA_ValidCase() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(0x00)
	s.Regs.BC.Set(addr)

	s.MockMMU.EXPECT().ReadByte(addr).Return(value, nil)

	result := s.factory.BCPtrToA()()
	s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
	s.Equal(value, s.Regs.A.Get())
}

func (s *unitTestSuite) TestBCPtrToA_InvalidRead() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(0x00)
	s.Regs.BC.Set(addr)

	s.MockMMU.EXPECT().ReadByte(addr).Return(value, testhelpers.ErrTesting1)

	result := s.factory.BCPtrToA()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestDEPtrToA_ValidCase() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(0x00)
	s.Regs.DE.Set(addr)

	s.MockMMU.EXPECT().ReadByte(addr).Return(value, nil)

	result := s.factory.DEPtrToA()()
	s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
	s.Equal(value, s.Regs.A.Get())
}

func (s *unitTestSuite) TestDEPtrToA_InvalidRead() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(0x00)
	s.Regs.DE.Set(addr)

	s.MockMMU.EXPECT().ReadByte(addr).Return(value, testhelpers.ErrTesting1)

	result := s.factory.DEPtrToA()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestHLIncToA_ValidCase() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(0x00)
	s.Regs.HL.Set(addr)

	s.MockMMU.EXPECT().ReadByte(addr).Return(value, nil)

	result := s.factory.HLIncToA()()
	s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
	s.Equal(value, s.Regs.A.Get())
	s.Equal(addr+1, s.Regs.HL.Get())
}

func (s *unitTestSuite) TestHLIncToA_InvalidRead() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(0x00)
	s.Regs.HL.Set(addr)

	s.MockMMU.EXPECT().ReadByte(addr).Return(value, testhelpers.ErrTesting1)

	result := s.factory.HLIncToA()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestHLDecToA_ValidCase() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(0x00)
	s.Regs.HL.Set(addr)

	s.MockMMU.EXPECT().ReadByte(addr).Return(value, nil)

	result := s.factory.HLDecToA()()
	s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
	s.Equal(value, s.Regs.A.Get())
	s.Equal(addr-1, s.Regs.HL.Get())
}

func (s *unitTestSuite) TestHLDecToA_InvalidRead() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(0x00)
	s.Regs.HL.Set(addr)

	s.MockMMU.EXPECT().ReadByte(addr).Return(value, testhelpers.ErrTesting1)

	result := s.factory.HLDecToA()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
