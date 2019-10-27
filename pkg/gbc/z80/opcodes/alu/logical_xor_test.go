package alu

import (
	"fmt"
	"nebula-go/pkg/common/testhelpers"
	"nebula-go/pkg/gbc/memory/registers"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

var _xorTestCases = []aOpTestCase{
	{0x00, 0x00, 0x00, z80lib.FlagsCleared, z80lib.ZF},
	{0x00, 0x00, 0x00, z80lib.FlagsFullSet, z80lib.ZF},
	{0x0F, 0x00, 0x0F, z80lib.FlagsCleared, z80lib.FlagsCleared},
	{0x0F, 0x0F, 0x00, z80lib.FlagsCleared, z80lib.ZF},
}

func (s *unitTestSuite) TestXorByteToA() {
	reg := registers.NewByte(0x00)
	fn := s.factory.XorByteToA(reg)

	for _, c := range _xorTestCases {
		name := fmt.Sprintf("xor a with register (a = %#v, reg = %#v)", c.initialValue, c.otherValue)

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

func (s *unitTestSuite) TestXorHLPtrToA_ValidCase() {
	fn := s.factory.XorHLPtrToA()

	for _, c := range _xorTestCases {
		name := fmt.Sprintf("xor a with hl ptr (a = %#v, hlPtr = %#v)", c.initialValue, c.otherValue)

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

func (s *unitTestSuite) TestXorHLPtrToA_FailedRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.XorHLPtrToA()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestXorD8ToA_ValidCase() {
	fn := s.factory.XorD8ToA()

	for _, c := range _xorTestCases {
		name := fmt.Sprintf("xor a with d8 (a = %#v, d8 = %#v)", c.initialValue, c.otherValue)

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

func (s *unitTestSuite) TestXorD8ToA_FailedRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.PC+1).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.XorD8ToA()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
