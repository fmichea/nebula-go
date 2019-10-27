package controlflow

import (
	"nebula-go/pkg/common/testhelpers"
	"nebula-go/pkg/gbc/memory/registers"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) callInterrupt(interrupt z80lib.Interrupt) {
	offsetSP := s.Regs.SP.Get() - 2
	pc := s.Regs.PC

	// Current PC should be different than the interrupt address.
	s.Assert().NotEqual(interrupt.Addr(), pc)

	// PC+1 is going to be pushed on the stack.
	s.MockMMU.EXPECT().WriteDByte(offsetSP, pc+1).Return(nil)

	// Build and call the interrupt caller
	callInterrupt := s.factory.CallInterrupt(interrupt)
	s.Require().NotNil(callInterrupt)

	// Result should have a size of 0 (since it is an opcode that modifies PC) and a clock of 16 ticks.
	result := callInterrupt()
	s.Equal(opcodeslib.OpcodeSuccess(1, 16), result)

	// Now PC is the interrupt's address, and SP has been offset.
	s.Assert().Equal(interrupt.Addr(), s.Regs.PC)
	s.Assert().Equal(offsetSP, s.Regs.SP.Get())
}

func (s *unitTestSuite) TestCallInterrupt_ValidCallRst00() {
	s.callInterrupt(z80lib.Rst00h)
}

func (s *unitTestSuite) TestCallInterrupt_ValidCallRst38() {
	s.callInterrupt(z80lib.Rst38h)
}

func (s *unitTestSuite) TestCallInterrupt_InvalidWrite() {
	offsetSP := s.Regs.SP.Get() - 2
	pc := s.Regs.PC

	s.MockMMU.EXPECT().WriteDByte(offsetSP, pc+1).Return(testhelpers.ErrTesting1)

	// Build and call the interrupt caller
	callInterrupt := s.factory.CallInterrupt(z80lib.Rst00h)
	s.Require().NotNil(callInterrupt)

	// This time the call returns the error raised by MMU.
	result := callInterrupt()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) callConditionalFuncValid(fn func(registers.Flag) opcodeslib.Opcode, value, shouldCall bool) {
	addr := uint16(0x1234)

	offsetSP := s.Regs.SP.Get() - 2
	pc := s.Regs.PC

	s.Regs.F.ZF.SetBool(value)

	// Current PC should be different than the call address.
	s.Assert().NotEqual(addr, pc)

	if shouldCall {
		// Address is going to be read right after opcode
		s.MockMMU.EXPECT().ReadDByte(pc+1).Return(addr, nil)

		// PC+3 is going to be pushed on the stack.
		s.MockMMU.EXPECT().WriteDByte(offsetSP, pc+3).Return(nil)
	}

	// Build and call the interrupt caller
	callConditional := fn(s.Regs.F.ZF)
	s.Require().NotNil(callConditional)

	result := callConditional()

	if shouldCall {
		s.Equal(opcodeslib.OpcodeSuccess(3, 24), result)
		s.Equal(addr, s.Regs.PC)
	} else {
		s.Equal(opcodeslib.OpcodeSuccess(3, 12), result)
		s.Equal(pc, s.Regs.PC)
	}
}

func (s *unitTestSuite) callConditionalAddressReadFailure(fn func(registers.Flag) opcodeslib.Opcode, value bool) {
	addr := uint16(0x1234)
	pc := s.Regs.PC

	s.Regs.F.ZF.SetBool(value)

	// Current PC should be different than the interrupt address.
	s.Assert().NotEqual(addr, pc)

	// Address is going to be read right after opcode
	s.MockMMU.EXPECT().ReadDByte(pc+1).Return(addr, testhelpers.ErrTesting1)

	// Build and call the interrupt caller
	callConditional := fn(s.Regs.F.ZF)
	s.Require().NotNil(callConditional)

	result := callConditional()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) callConditionalAddressPushFailure(fn func(registers.Flag) opcodeslib.Opcode, value bool) {
	addr := uint16(0x1234)

	offsetSP := s.Regs.SP.Get() - 2
	pc := s.Regs.PC

	s.Regs.F.ZF.SetBool(value)

	// Current PC should be different than the interrupt address.
	s.Assert().NotEqual(addr, pc)

	// Address is going to be read right after opcode
	s.MockMMU.EXPECT().ReadDByte(pc+1).Return(addr, nil)

	// Write to stack failure.
	s.MockMMU.EXPECT().WriteDByte(offsetSP, pc+3).Return(testhelpers.ErrTesting1)

	// Build and call the interrupt caller
	callConditional := fn(s.Regs.F.ZF)
	s.Require().NotNil(callConditional)

	result := callConditional()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestCallIf_FlagIsFalse() {
	s.callConditionalFuncValid(s.factory.CallIf, false, false)
}

func (s *unitTestSuite) TestCallIf_FlagIsTrue() {
	s.callConditionalFuncValid(s.factory.CallIf, true, true)

}

func (s *unitTestSuite) TestCallIf_InvalidAddressRead() {
	s.callConditionalAddressReadFailure(s.factory.CallIf, true)
}

func (s *unitTestSuite) TestCallIf_AddressPushFailure() {
	s.callConditionalAddressPushFailure(s.factory.CallIf, true)
}

func (s *unitTestSuite) TestCallIfNot_FlagIsFalse() {
	s.callConditionalFuncValid(s.factory.CallIfNot, false, true)
}

func (s *unitTestSuite) TestCallIfNot_FlagIsTrue() {
	s.callConditionalFuncValid(s.factory.CallIfNot, true, false)
}

func (s *unitTestSuite) TestCallIfNot_InvalidAddressRead() {
	s.callConditionalAddressReadFailure(s.factory.CallIfNot, false)
}

func (s *unitTestSuite) TestCallIfNot_AddressPushFailure() {
	s.callConditionalAddressPushFailure(s.factory.CallIfNot, false)
}

func (s *unitTestSuite) TestCall_ValidCase() {
	addr := uint16(0x1234)

	offsetSP := s.Regs.SP.Get() - 2
	pc := s.Regs.PC

	// Current PC should be different than the interrupt address.
	s.Assert().NotEqual(addr, pc)

	// Address is going to be read right after opcode
	s.MockMMU.EXPECT().ReadDByte(pc+1).Return(addr, nil)

	// PC+3 is going to be pushed on the stack.
	s.MockMMU.EXPECT().WriteDByte(offsetSP, pc+3).Return(nil)

	// Build and call the interrupt caller
	call := s.factory.Call()
	s.Require().NotNil(call)

	result := call()

	s.Equal(opcodeslib.OpcodeSuccess(3, 24), result)
	s.Equal(addr, s.Regs.PC)
}

func (s *unitTestSuite) TestCall_InvalidAddressRead() {
	addr := uint16(0x1234)

	pc := s.Regs.PC

	// Current PC should be different than the interrupt address.
	s.Assert().NotEqual(addr, pc)

	// Address is going to be read right after opcode
	s.MockMMU.EXPECT().ReadDByte(pc+1).Return(addr, testhelpers.ErrTesting1)

	// Build and call the interrupt caller
	call := s.factory.Call()
	s.Require().NotNil(call)

	result := call()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestCall_AddressPushFailure() {
	addr := uint16(0x1234)

	offsetSP := s.Regs.SP.Get() - 2
	pc := s.Regs.PC

	// Current PC should be different than the interrupt address.
	s.Assert().NotEqual(addr, pc)

	// Address is going to be read right after opcode
	s.MockMMU.EXPECT().ReadDByte(pc+1).Return(addr, nil)

	// PC+3 is going to be pushed on the stack.
	s.MockMMU.EXPECT().WriteDByte(offsetSP, pc+3).Return(testhelpers.ErrTesting1)

	// Build and call the interrupt caller
	call := s.factory.Call()
	s.Require().NotNil(call)

	result := call()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
