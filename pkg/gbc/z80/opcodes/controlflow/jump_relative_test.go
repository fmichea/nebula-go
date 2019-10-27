package controlflow

import (
	"nebula-go/pkg/common/testhelpers"
	opcodes_lib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) jumpRelativeFuncValid(fn func() opcodes_lib.OpcodeResult, shouldJump bool, offset uint8) {
	pc := s.Regs.PC
	addr := opcodes_lib.AddRelativeConst(pc, offset)

	if shouldJump {
		s.MockMMU.EXPECT().ReadByte(pc+1).Return(offset, nil)
	}

	result := fn()

	if shouldJump {
		s.Equal(opcodes_lib.OpcodeSuccess(2, 12), result)
		s.Equal(addr, s.Regs.PC)
	} else {
		s.Equal(opcodes_lib.OpcodeSuccess(2, 8), result)
		s.Equal(pc, s.Regs.PC)
	}
}

func (s *unitTestSuite) jumpRelativeFuncFailedRead(fn func() opcodes_lib.OpcodeResult) {
	pc := s.Regs.PC

	s.MockMMU.EXPECT().ReadByte(pc+1).Return(uint8(0), testhelpers.ErrTesting1)

	result := fn()
	s.Equal(opcodes_lib.OpcodeError(testhelpers.ErrTesting1), result)
	s.Equal(pc, s.Regs.PC)
}

func (s *unitTestSuite) jumpRelativeFuncValidNeg(fn func() opcodes_lib.OpcodeResult, shouldJump bool) {
	pc := s.Regs.PC

	// 0xF6 = -10
	s.jumpRelativeFuncValid(fn, shouldJump, 0xF6)
	s.True(!shouldJump || s.Regs.PC < pc)
}

func (s *unitTestSuite) jumpRelativeFuncValidPos(fn func() opcodes_lib.OpcodeResult, shouldJump bool) {
	pc := s.Regs.PC

	s.jumpRelativeFuncValid(fn, shouldJump, 0x0a)
	s.True(!shouldJump || pc < s.Regs.PC)
}

func (s *unitTestSuite) TestJumpRelative_ValidCase_Negative() {
	s.jumpRelativeFuncValidNeg(s.factory.JumpRelative(), true)
}

func (s *unitTestSuite) TestJumpRelative_ValidCase_Positive() {
	s.jumpRelativeFuncValidPos(s.factory.JumpRelative(), true)
}

func (s *unitTestSuite) TestJumpRelative_FailedRead() {
	s.jumpRelativeFuncFailedRead(s.factory.JumpRelative())
}

func (s *unitTestSuite) TestJumpRelativeIf_FlagIsTrue_Negative() {
	flag := s.Regs.F.ZF
	flag.SetBool(true)

	s.jumpRelativeFuncValidNeg(s.factory.JumpRelativeIf(flag), true)
}

func (s *unitTestSuite) TestJumpRelativeIf_FlagIsTrue_Positive() {
	flag := s.Regs.F.ZF
	flag.SetBool(true)

	s.jumpRelativeFuncValidPos(s.factory.JumpRelativeIf(flag), true)
}

func (s *unitTestSuite) TestJumpRelativeIf_FlagIsFalse_Negative() {
	flag := s.Regs.F.ZF
	flag.SetBool(false)

	s.jumpRelativeFuncValidNeg(s.factory.JumpRelativeIf(flag), false)
}

func (s *unitTestSuite) TestJumpRelativeIf_FlagIsFalse_Positive() {
	flag := s.Regs.F.ZF
	flag.SetBool(false)

	s.jumpRelativeFuncValidPos(s.factory.JumpRelativeIf(flag), false)
}

func (s *unitTestSuite) TestJumpRelativeIf_FailedRead() {
	flag := s.Regs.F.ZF
	flag.SetBool(true)

	s.jumpRelativeFuncFailedRead(s.factory.JumpRelativeIf(flag))
}

func (s *unitTestSuite) TestJumpRelativeIfNot_FlagIsTrue_Negative() {
	flag := s.Regs.F.ZF
	flag.SetBool(true)

	s.jumpRelativeFuncValidNeg(s.factory.JumpRelativeIfNot(flag), false)
}

func (s *unitTestSuite) TestJumpRelativeIfNot_FlagIsTrue_Positive() {
	flag := s.Regs.F.ZF
	flag.SetBool(true)

	s.jumpRelativeFuncValidPos(s.factory.JumpRelativeIfNot(flag), false)
}

func (s *unitTestSuite) TestJumpRelativeIfNot_FlagIsFalse_Negative() {
	flag := s.Regs.F.ZF
	flag.SetBool(false)

	s.jumpRelativeFuncValidNeg(s.factory.JumpRelativeIfNot(flag), true)
}

func (s *unitTestSuite) TestJumpRelativeIfNot_FlagIsFalse_Positive() {
	flag := s.Regs.F.ZF
	flag.SetBool(false)

	s.jumpRelativeFuncValidPos(s.factory.JumpRelativeIfNot(flag), true)
}

func (s *unitTestSuite) TestJumpRelativeIfNot_FailedRead() {
	flag := s.Regs.F.ZF
	flag.SetBool(false)

	s.jumpRelativeFuncFailedRead(s.factory.JumpRelativeIfNot(flag))
}
