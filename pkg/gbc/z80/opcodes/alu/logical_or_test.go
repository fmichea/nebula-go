package alu

import (
	"fmt"

	"nebula-go/pkg/common/testhelpers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

var _orTestCases = []aOpTestCase{
	{0x00, 0x00, 0x00, registers.FlagsCleared, registers.ZF},
	{0x00, 0x00, 0x00, registers.FlagsFullSet, registers.ZF},
	{0x0F, 0xF0, 0xFF, registers.FlagsCleared, registers.FlagsCleared},
	{0x0F, 0xF0, 0xFF, registers.FlagsFullSet, registers.FlagsCleared},
	{0xAA, 0x55, 0xFF, registers.FlagsCleared, registers.FlagsCleared},
	{0xC2, 0x22, 0xE2, registers.FlagsFullSet, registers.FlagsCleared},
}

func (s *unitTestSuite) TestOrByteToA() {
	reg := registerslib.NewByte(0x00)
	fn := s.factory.OrByteToA(reg)

	for _, c := range _orTestCases {
		name := fmt.Sprintf("or a with register (a = %#v, reg = %#v)", c.initialValue, c.otherValue)

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

func (s *unitTestSuite) TestOrHLPtrToA_ValidCase() {
	fn := s.factory.OrHLPtrToA()

	for _, c := range _orTestCases {
		name := fmt.Sprintf("or a with hl ptr (a = %#v, hlPtr = %#v)", c.initialValue, c.otherValue)

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

func (s *unitTestSuite) TestOrHLPtrToA_FailedRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.OrHLPtrToA()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestOrD8ToA_ValidCase() {
	fn := s.factory.OrD8ToA()

	for _, c := range _orTestCases {
		name := fmt.Sprintf("or a with d8 (a = %#v, d8 = %#v)", c.initialValue, c.otherValue)

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

func (s *unitTestSuite) TestOrD8ToA_FailedRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.PC+1).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.OrD8ToA()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
