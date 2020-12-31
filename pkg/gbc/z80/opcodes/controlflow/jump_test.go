package controlflow

import (
	"nebula-go/pkg/common/testhelpers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) jumpFuncValid(fn func() opcodeslib.OpcodeResult, shouldJump bool) {
	addr := uint16(0x1234)

	pc := s.Regs.PC

	// Current PC should be different than the jump address
	s.NotEqual(addr, pc)

	if shouldJump {
		// Address will be read right after opcode.
		s.MockMMU.EXPECT().ReadDByte(pc+1).Return(addr, nil)
	}

	result := fn()

	if shouldJump {
		s.Equal(opcodeslib.OpcodeSuccess(0, 16), result)
		s.Equal(addr, s.Regs.PC)
	} else {
		s.Equal(opcodeslib.OpcodeSuccess(3, 12), result)
		s.Equal(pc, s.Regs.PC)
	}
}

func (s *unitTestSuite) jumpFuncInvalidRead(fn func() opcodeslib.OpcodeResult) {
	addr := uint16(0x1234)

	pc := s.Regs.PC

	// Current PC should be different than the jump address
	s.NotEqual(addr, pc)

	// Address will be read right after opcode.
	s.MockMMU.EXPECT().ReadDByte(pc+1).Return(addr, testhelpers.ErrTesting1)

	result := fn()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestJumpIf_FlagIsTrue() {
	flag := s.Regs.F.ZF
	flag.SetBool(true)

	s.jumpFuncValid(s.factory.JumpIf(flag), true)
}

func (s *unitTestSuite) TestJumpIf_FlagIsFalse() {
	flag := s.Regs.F.ZF
	flag.SetBool(false)

	s.jumpFuncValid(s.factory.JumpIf(flag), false)
}

func (s *unitTestSuite) TestJumpIf_InvalidRead() {
	flag := s.Regs.F.ZF
	flag.SetBool(true)

	s.jumpFuncInvalidRead(s.factory.JumpIf(flag))
}

func (s *unitTestSuite) TestJumpIfNot_FlagIsTrue() {
	flag := s.Regs.F.ZF
	flag.SetBool(true)

	s.jumpFuncValid(s.factory.JumpIfNot(flag), false)
}

func (s *unitTestSuite) TestJumpIfNot_FlagIsFalse() {
	flag := s.Regs.F.ZF
	flag.SetBool(false)

	s.jumpFuncValid(s.factory.JumpIfNot(flag), true)
}

func (s *unitTestSuite) TestJumpIfNot_InvalidRead() {
	flag := s.Regs.F.ZF
	flag.SetBool(false)

	s.jumpFuncInvalidRead(s.factory.JumpIfNot(flag))
}

func (s *unitTestSuite) TestJump_ValidCase() {
	s.jumpFuncValid(s.factory.Jump(), true)
}

func (s *unitTestSuite) TestJump_InvalidRead() {
	s.jumpFuncInvalidRead(s.factory.Jump())
}

func (s *unitTestSuite) TestJumpHL_ValidCase() {
	addr := uint16(0x1234)

	s.Regs.HL.Set(addr)

	fn := s.factory.JumpHL()
	result := fn()

	s.Equal(opcodeslib.OpcodeSuccess(0, 4), result)
	s.Equal(addr, s.Regs.PC)
}
