package alu

import (
	"fmt"

	"nebula-go/pkg/common/testhelpers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

var _adcTestCases = []aOpTestCase{
	// Nothing to add, result is 0.
	{0x00, 0x00, 0x00, registers.FlagsCleared, registers.ZF},
	// Carry is accounted for.
	{0x00, 0x00, 0x01, registers.FlagsFullSet, registers.FlagsCleared},
	// Half-carry overflow.
	{0x0F, 0x0F, 0x1E, registers.FlagsCleared, registers.HC},
	{0x0F, 0x0F, 0x1F, registers.FlagsFullSet, registers.HC},
	// Carry
	{0xFF, 0xF, 0x0E, registers.FlagsCleared, registers.HC | registers.CY},
}

func (s *unitTestSuite) TestAdcByteToA() {
	reg := registerslib.NewByte(0x00)
	fn := s.factory.AdcByteToA(reg)

	for _, c := range _adcTestCases {
		name := fmt.Sprintf("adc a with register (a = %#v, reg = %#v)", c.initialValue, c.otherValue)

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

func (s *unitTestSuite) TestAdcHLPtrToA_ValidCase() {
	fn := s.factory.AdcHLPtrToA()

	for _, c := range _adcTestCases {
		name := fmt.Sprintf("adc a with hl ptr (a = %#v, hlPtr = %#v)", c.initialValue, c.otherValue)

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

func (s *unitTestSuite) TestAdcHLPtrToA_FailedRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.AdcHLPtrToA()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestAdcD8ToA_ValidCase() {
	fn := s.factory.AdcD8ToA()

	for _, c := range _adcTestCases {
		name := fmt.Sprintf("adc a with d8 (a = %#v, d8 = %#v)", c.initialValue, c.otherValue)

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

func (s *unitTestSuite) TestAdcD8ToA_FailedRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.PC+1).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.AdcD8ToA()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
