package alu

import (
	"fmt"

	"nebula-go/pkg/common/testhelpers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

var _addTestCases = []aOpTestCase{
	// Nothing to add, result is 0.
	{0x00, 0x00, 0x00, registers.FlagsCleared, registers.ZF},
	{0x00, 0x00, 0x00, registers.FlagsFullSet, registers.ZF},
	// Small add.
	{0x00, 0x01, 0x01, registers.FlagsFullSet, registers.FlagsCleared},
	// Half-carry overflow.
	{0x0F, 0x0F, 0x1E, registers.FlagsCleared, registers.HC},
	{0x0F, 0x0F, 0x1E, registers.FlagsFullSet, registers.HC},
	// Carry
	{0xFF, 0xF, 0x0E, registers.FlagsCleared, registers.HC | registers.CY},
}

func (s *unitTestSuite) TestAddByteToA() {
	reg := registerslib.NewByte(0x00)
	fn := s.factory.AddByteToA(reg)

	for _, c := range _addTestCases {
		name := fmt.Sprintf("add a with register (a = %#v, reg = %#v)", c.initialValue, c.otherValue)

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

func (s *unitTestSuite) TestAddHLPtrToA_ValidCase() {
	fn := s.factory.AddHLPtrToA()

	for _, c := range _addTestCases {
		name := fmt.Sprintf("add a with hl ptr (a = %#v, hlPtr = %#v)", c.initialValue, c.otherValue)

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

func (s *unitTestSuite) TestAddHLPtrToA_FailedRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.AddHLPtrToA()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestAddD8ToA_ValidCase() {
	fn := s.factory.AddD8ToA()

	for _, c := range _addTestCases {
		name := fmt.Sprintf("add a with d8 (a = %#v, d8 = %#v)", c.initialValue, c.otherValue)

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

func (s *unitTestSuite) TestAddD8ToA_FailedRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.PC+1).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.AddD8ToA()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestAddDByteToHL() {
	cases := []struct {
		initialValue uint16
		otherValue   uint16
		resultValue  uint16

		initialFlags uint8
		resultFlags  uint8
	}{
		// Zero sum results in zero result (zero-flag not affected by this op)
		{0x0000, 0x0000, 0x0000, registers.FlagsCleared, registers.FlagsCleared},
		{0x0000, 0x0000, 0x0000, registers.FlagsFullSet, registers.ZF},
		// Some simple cases, no carry. (again zero-flag not affected by this op)
		{0x0000, 0x0000, 0x0000, registers.FlagsCleared, registers.FlagsCleared},
		{0x0008, 0x0008, 0x0010, registers.FlagsFullSet, registers.ZF},
		// Half carry does not happen on bit 9, but on bit 13 actually.
		{0x0080, 0x0080, 0x0100, registers.FlagsFullSet, registers.ZF},
		{0x0200, 0x0F00, 0x1100, registers.FlagsCleared, registers.HC},
		// Carry works as usual.
		{0xFFFF, 0x000F, 0x000E, registers.FlagsCleared, registers.CY | registers.HC},
	}

	reg := registerslib.NewDByte(0x00)
	fn := s.factory.AddDByteToHL(reg)

	for _, c := range cases {
		name := fmt.Sprintf("add dbyte reg to hl (hl = %#v, reg = %#v)", c.initialValue, c.otherValue)

		s.Run(name, func() {
			reg.Set(c.otherValue)
			s.Regs.HL.Set(c.initialValue)
			s.Regs.F.Set(c.initialFlags)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
			s.Equal(c.resultValue, s.Regs.HL.Get())
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestAddR8ToSP_ValidCase() {
	cases := []struct {
		initialValue uint16
		r8           uint8
		resultValue  uint16

		initialFlags uint8
		resultFlags  uint8
	}{
		{0xABCD, 0x00, 0xABCD, registers.FlagsFullSet, registers.FlagsCleared},

		{0xABCD, 0x01, 0xABCE, registers.FlagsCleared, registers.FlagsCleared},
		{0xABCD, 0x01, 0xABCE, registers.FlagsFullSet, registers.FlagsCleared},
		{0xABCD, 0x04, 0xABD1, registers.FlagsCleared, registers.HC},
		{0xABFD, 0x10, 0xAC0D, registers.FlagsCleared, registers.CY},
		{0xABCD, 0x73, 0xAC40, registers.FlagsCleared, registers.HC | registers.CY},

		// Addition with negative:
		{0xAB1D, 0x85, 0xAAA2, registers.FlagsCleared, registers.HC},
		{0xABF0, 0xCE, 0xABBE, registers.FlagsCleared, registers.CY},
		{0xABF5, 0xCE, 0xABC3, registers.FlagsCleared, registers.HC | registers.CY},
	}

	fn := s.factory.AddR8ToSP()

	for _, c := range cases {
		name := fmt.Sprintf("add signed int8 to sp (sp = %#v, r8 = %#v)", c.initialValue, c.r8)

		s.Run(name, func() {
			pc := s.Regs.PC

			s.MockMMU.EXPECT().ReadByte(pc+1).Return(c.r8, nil)

			s.Regs.F.Set(c.initialFlags)
			s.Regs.SP.Set(c.initialValue)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(2, 16), result)
			s.Equal(c.resultValue, s.Regs.SP.Get())
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestAddR8ToSP_InvalidRead() {
	pc := s.Regs.PC

	s.MockMMU.EXPECT().ReadByte(pc+1).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.AddR8ToSP()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
