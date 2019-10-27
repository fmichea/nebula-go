package load

import (
	"nebula-go/pkg/common/testhelpers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) TestAToAddress_ValidCase() {
	a := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(a)

	s.MockMMU.EXPECT().ReadDByte(s.Regs.PC+1).Return(addr, nil)
	s.MockMMU.EXPECT().WriteByte(addr, a).Return(nil)

	result := s.factory.AToAddress()()
	s.Equal(opcodeslib.OpcodeSuccess(3, 16), result)
}

func (s *unitTestSuite) TestAToAddress_InvalidRead() {
	a := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(a)

	s.MockMMU.EXPECT().ReadDByte(s.Regs.PC+1).Return(addr, testhelpers.ErrTesting1)

	result := s.factory.AToAddress()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestAToAddress_InvalidWrite() {
	a := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(a)

	s.MockMMU.EXPECT().ReadDByte(s.Regs.PC+1).Return(addr, nil)
	s.MockMMU.EXPECT().WriteByte(addr, a).Return(testhelpers.ErrTesting1)

	result := s.factory.AToAddress()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestAToBCPtr_ValidCase() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(value)
	s.Regs.BC.Set(addr)

	s.MockMMU.EXPECT().WriteByte(addr, value).Return(nil)

	result := s.factory.AToBCPtr()()
	s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
}

func (s *unitTestSuite) TestAToBCPtr_InvalidRead() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(value)
	s.Regs.BC.Set(addr)

	s.MockMMU.EXPECT().WriteByte(addr, value).Return(testhelpers.ErrTesting1)

	result := s.factory.AToBCPtr()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestAToDEPtr_ValidCase() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(value)
	s.Regs.DE.Set(addr)

	s.MockMMU.EXPECT().WriteByte(addr, value).Return(nil)

	result := s.factory.AToDEPtr()()
	s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
}

func (s *unitTestSuite) TestAToDEPtr_InvalidRead() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(value)
	s.Regs.DE.Set(addr)

	s.MockMMU.EXPECT().WriteByte(addr, value).Return(testhelpers.ErrTesting1)

	result := s.factory.AToDEPtr()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestAToHLInc_ValidCase() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(value)
	s.Regs.HL.Set(addr)

	s.MockMMU.EXPECT().WriteByte(addr, value).Return(nil)

	result := s.factory.AToHLInc()()
	s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
	s.Equal(addr+1, s.Regs.HL.Get())
}

func (s *unitTestSuite) TestAToHLInc_InvalidRead() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(value)
	s.Regs.HL.Set(addr)

	s.MockMMU.EXPECT().WriteByte(addr, value).Return(testhelpers.ErrTesting1)

	result := s.factory.AToHLInc()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestAToHLDec_ValidCase() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(value)
	s.Regs.HL.Set(addr)

	s.MockMMU.EXPECT().WriteByte(addr, value).Return(nil)

	result := s.factory.AToHLDec()()
	s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
	s.Equal(addr-1, s.Regs.HL.Get())
}

func (s *unitTestSuite) TestAToHLDec_InvalidRead() {
	value := uint8(0xDE)
	addr := uint16(0xABCD)

	s.Regs.A.Set(value)
	s.Regs.HL.Set(addr)

	s.MockMMU.EXPECT().WriteByte(addr, value).Return(testhelpers.ErrTesting1)

	result := s.factory.AToHLDec()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
