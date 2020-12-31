package alu

import (
	"fmt"

	"nebula-go/pkg/common/testhelpers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

// These tests are the same as substract tests, but the the result value is never changed.
var _compareTestCases = []aOpTestCase{
	// Nothing substracted, result is 0.
	{0x00, 0x00, 0x00, registers.FlagsCleared, registers.ZF | registers.NE},
	{0xFF, 0xFF, 0xFF, registers.FlagsFullSet, registers.ZF | registers.NE},
	// Value and carry substract to 0.
	{0x03, 0x02, 0x03, registers.FlagsFullSet, registers.NE},
	// Carry borrows only on bit 4, no CY in result.
	{0x22, 0x03, 0x22, registers.FlagsCleared, registers.NE | registers.HC},
	// Carry borrows only on bit 4, no CY in result.
	{0x22, 0x03, 0x22, registers.FlagsFullSet, registers.NE | registers.HC},
	// Carry borrows only on bit 4, no CY in result, case were bit 4 is not set in result.
	{0x82, 0x03, 0x82, registers.FlagsCleared, registers.NE | registers.HC},
	// Complete borrow.
	{0x00, 0x0F, 0x00, registers.FlagsCleared, registers.NE | registers.HC | registers.CY},
	// No borrow, not zero.
	{0x02, 0x00, 0x02, registers.FlagsFullSet, registers.NE},
	// No borrow, not zero.
	{0x02, 0x01, 0x02, registers.FlagsCleared, registers.NE},
}

func (s *unitTestSuite) TestCompareByteToA() {
	reg := registerslib.NewByte(0x00)
	fn := s.factory.CompareByteToA(reg)

	for _, c := range _compareTestCases {
		name := fmt.Sprintf("compare a with register (a = %#v, reg = %#v)", c.initialValue, c.otherValue)

		s.Run(name, func() {
			reg.Set(c.otherValue)
			s.Regs.A.Set(c.initialValue)
			s.Regs.F.Set(c.initialFlags)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
			s.Equal(c.resultValue, s.Regs.A.Get())
			s.EqualFlags(c.resultFlags)
			s.Equal(c.otherValue, reg.Get())
		})
	}
}

func (s *unitTestSuite) TestCompareHLPtrToA_ValidCase() {
	fn := s.factory.CompareHLPtrToA()

	for _, c := range _compareTestCases {
		name := fmt.Sprintf("compare a with hl ptr (a = %#v, hlPtr = %#v)", c.initialValue, c.otherValue)

		s.Run(name, func() {
			s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(c.otherValue, nil)

			s.Regs.A.Set(c.initialValue)
			s.Regs.F.Set(c.initialFlags)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
			s.Equal(c.resultValue, s.Regs.A.Get())
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestCompareHLPtrToA_FailedRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.CompareHLPtrToA()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestCompareD8ToA_ValidCase() {
	fn := s.factory.CompareD8ToA()

	for _, c := range _compareTestCases {
		name := fmt.Sprintf("compare a with d8 (a = %#v, d8 = %#v)", c.initialValue, c.otherValue)

		s.Run(name, func() {
			s.MockMMU.EXPECT().ReadByte(s.Regs.PC+1).Return(c.otherValue, nil)

			s.Regs.A.Set(c.initialValue)
			s.Regs.F.Set(c.initialFlags)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(2, 8), result)
			s.Equal(c.resultValue, s.Regs.A.Get())
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestCompareD8ToA_FailedRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.PC+1).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.CompareD8ToA()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
