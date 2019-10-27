package alu

import (
	"fmt"
	"nebula-go/pkg/common/testhelpers"
	"nebula-go/pkg/gbc/memory/registers"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

var _subTestCases = []aOpTestCase{
	// Nothing substracted, result is 0.
	{0x00, 0x00, 0x00, z80lib.FlagsCleared, z80lib.ZF | z80lib.NE},
	// Carry borrows, both cy and hc get set.
	{0x00, 0x01, 0xFF, z80lib.FlagsFullSet, z80lib.NE | z80lib.CY | z80lib.HC},
	// Carry borrows only on bit 4, no CY in result.
	{0x22, 0x03, 0x1F, z80lib.FlagsCleared, z80lib.NE | z80lib.HC},
	// Carry borrows only on bit 4, no CY in result.
	{0x22, 0x03, 0x1F, z80lib.FlagsFullSet, z80lib.NE | z80lib.HC},
	// Carry borrows only on bit 4, no CY in result, case were bit 4 is not set in result.
	{0x82, 0x03, 0x7F, z80lib.FlagsCleared, z80lib.NE | z80lib.HC},
	// No borrow, not zero.
	{0x02, 0x01, 0x01, z80lib.FlagsCleared, z80lib.NE},
}

func (s *unitTestSuite) TestSubByteToA() {
	reg := registers.NewByte(0x00)
	fn := s.factory.SubByteToA(reg)

	for _, c := range _subTestCases {
		name := fmt.Sprintf("sub a with register (a = %#v, reg = %#v)", c.initialValue, c.otherValue)

		s.Run(name, func() {
			reg.Set(c.otherValue)
			s.Regs.A.Set(c.initialValue)
			s.Regs.F.Set(c.initialFlags)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
			s.Equal(c.resultValue, s.Regs.A.Get())
			s.Equal(c.otherValue, reg.Get())
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestSubHLPtrToA_ValidCase() {
	fn := s.factory.SubHLPtrToA()

	for _, c := range _subTestCases {
		name := fmt.Sprintf("sub a with hl ptr (a = %#v, hlPtr = %#v)", c.initialValue, c.otherValue)

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

func (s *unitTestSuite) TestSubHLPtrToA_FailedRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.SubHLPtrToA()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestSubD8ToA_ValidCase() {
	fn := s.factory.SubD8ToA()

	for _, c := range _subTestCases {
		name := fmt.Sprintf("sub a with d8 (a = %#v, d8 = %#v)", c.initialValue, c.otherValue)

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

func (s *unitTestSuite) TestSubD8ToA_FailedRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.PC+1).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.SubD8ToA()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
