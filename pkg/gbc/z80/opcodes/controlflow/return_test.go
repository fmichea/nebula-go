package controlflow

import (
	"nebula-go/pkg/common/testhelpers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) returnIfFuncValid(fn func() opcodeslib.OpcodeResult, shouldReturn bool) {
	pc := s.Regs.PC
	sp := s.Regs.SP.Get()
	addr := uint16(0x1234)

	if shouldReturn {
		s.MockMMU.EXPECT().ReadDByte(sp).Return(addr, nil)
	}

	result := fn()

	if shouldReturn {
		s.Equal(opcodeslib.OpcodeSuccess(0, 20), result)
		s.Equal(addr, s.Regs.PC)
		s.Equal(sp+2, s.Regs.SP.Get())
	} else {
		s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
		s.Equal(pc, s.Regs.PC)
		s.Equal(sp, s.Regs.SP.Get())
	}
}

func (s *unitTestSuite) returnIfFailedAddressRead(fn func() opcodeslib.OpcodeResult) {
	s.MockMMU.EXPECT().ReadDByte(s.Regs.SP.Get()).Return(uint16(0), testhelpers.ErrTesting1)

	result := fn()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) returnFuncValid(fn func() opcodeslib.OpcodeResult) {
	pc := s.Regs.PC
	sp := s.Regs.SP.Get()
	addr := uint16(0x1234)

	s.NotEqual(addr, pc)

	s.MockMMU.EXPECT().ReadDByte(sp).Return(addr, nil)

	result := fn()
	s.Equal(opcodeslib.OpcodeSuccess(0, 16), result)
	s.Equal(sp+2, s.Regs.SP.Get())
	s.Equal(addr, s.Regs.PC)
}

func (s *unitTestSuite) returnFuncFailedAddressRead(fn func() opcodeslib.OpcodeResult) {
	s.MockMMU.EXPECT().ReadDByte(s.Regs.SP.Get()).Return(uint16(0), testhelpers.ErrTesting1)
	result := fn()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestReturnIf_FlagIsTrue() {
	flag := s.Regs.F.ZF
	flag.SetBool(true)

	s.returnIfFuncValid(s.factory.ReturnIf(flag), true)
}

func (s *unitTestSuite) TestReturnIf_FlagIsFalse() {
	flag := s.Regs.F.ZF
	flag.SetBool(false)

	s.returnIfFuncValid(s.factory.ReturnIf(flag), false)
}

func (s *unitTestSuite) TestReturnIf_FailedRead() {
	flag := s.Regs.F.ZF
	flag.SetBool(true)

	s.returnIfFailedAddressRead(s.factory.ReturnIf(flag))
}

func (s *unitTestSuite) TestReturnIfNot_FlagIsTrue() {
	flag := s.Regs.F.ZF
	flag.SetBool(true)

	s.returnIfFuncValid(s.factory.ReturnIfNot(flag), false)
}

func (s *unitTestSuite) TestReturnIfNot_FlagIsFalse() {
	flag := s.Regs.F.ZF
	flag.SetBool(false)

	s.returnIfFuncValid(s.factory.ReturnIfNot(flag), true)
}

func (s *unitTestSuite) TestReturnIfNot_FailedRead() {
	flag := s.Regs.F.ZF
	flag.SetBool(false)

	s.returnIfFailedAddressRead(s.factory.ReturnIfNot(flag))
}

func (s *unitTestSuite) TestReturn_ValidCase() {
	s.returnFuncValid(s.factory.Return())
}

func (s *unitTestSuite) TestReturn_FailedRead() {
	s.returnFuncFailedAddressRead(s.factory.Return())
}

func (s *unitTestSuite) TestReturnInterrupt_ValidCase_IMEFalse() {
	s.Regs.IME = false

	s.returnFuncValid(s.factory.ReturnInterrupt())
	s.True(s.Regs.IME)
}

func (s *unitTestSuite) TestReturnInterrupt_ValidCase_IMETrue() {
	s.Regs.IME = true

	s.returnFuncValid(s.factory.ReturnInterrupt())
	s.True(s.Regs.IME)
}

func (s *unitTestSuite) TestReturnInterrupt_FailedRead() {
	s.returnFuncFailedAddressRead(s.factory.ReturnInterrupt())
}
