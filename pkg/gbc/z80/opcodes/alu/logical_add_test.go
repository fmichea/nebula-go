package alu

import (
	"fmt"
	"nebula-go/pkg/common/testhelpers"
	"nebula-go/pkg/gbc/memory/registers"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

var _addTestCases = []aOpTestCase{
	// Nothing to add, result is 0.
	{0x00, 0x00, 0x00, z80lib.FlagsCleared, z80lib.ZF},
	{0x00, 0x00, 0x00, z80lib.FlagsFullSet, z80lib.ZF},
	// Small add.
	{0x00, 0x01, 0x01, z80lib.FlagsFullSet, z80lib.FlagsCleared},
	// Half-carry overflow.
	{0x0F, 0x0F, 0x1E, z80lib.FlagsCleared, z80lib.HC},
	{0x0F, 0x0F, 0x1E, z80lib.FlagsFullSet, z80lib.HC},
	// Carry
	{0xFF, 0xF, 0x0E, z80lib.FlagsCleared, z80lib.HC | z80lib.CY},
}

func (s *unitTestSuite) TestAddByteToA() {
	reg := registers.NewByte(0x00)
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
		{0x0000, 0x0000, 0x0000, z80lib.FlagsCleared, z80lib.FlagsCleared},
		{0x0000, 0x0000, 0x0000, z80lib.FlagsFullSet, z80lib.ZF},
		// Some simple cases, no carry. (again zero-flag not affected by this op)
		{0x0000, 0x0000, 0x0000, z80lib.FlagsCleared, z80lib.FlagsCleared},
		{0x0008, 0x0008, 0x0010, z80lib.FlagsFullSet, z80lib.ZF},
		// Half carry does not happen on bit 9, but on bit 13 actually.
		{0x0080, 0x0080, 0x0100, z80lib.FlagsFullSet, z80lib.ZF},
		{0x0200, 0x0F00, 0x1100, z80lib.FlagsCleared, z80lib.HC},
		// Carry works as usual.
		{0xFFFF, 0x000F, 0x000E, z80lib.FlagsCleared, z80lib.CY | z80lib.HC},
	}

	reg := registers.NewDByte(0x00)
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
		otherValue   uint8
		resultValue  uint16
	}{
		// Simple adding value.
		{0xFFF0, 0x00, 0xFFF0},
		{0xFFF0, 0x04, 0xFFF4},
		// Negative value.
		{0xFFF0, 0xF6, 0xFFE6},
	}

	fn := s.factory.AddR8ToSP()

	for _, c := range cases {
		name := fmt.Sprintf("add signed int8 to sp (sp = %#v, r8 = %#v)", c.initialValue, c.otherValue)

		s.Run(name, func() {
			pc := s.Regs.PC

			s.MockMMU.EXPECT().ReadByte(pc+1).Return(c.otherValue, nil)

			s.Regs.SP.Set(c.initialValue)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(2, 16), result)
			s.Equal(c.resultValue, s.Regs.SP.Get())
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
