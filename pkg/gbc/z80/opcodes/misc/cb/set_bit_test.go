package cb

import (
	"fmt"

	"nebula-go/pkg/common/testhelpers"
	lib2 "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

var _setBitCases = []cbBitOpTestCase{
	{0xF0, 0, 0xF1, registers.FlagsCleared, registers.FlagsCleared},
	{0xF0, 0, 0xF1, registers.FlagsFullSet, registers.FlagsFullSet},

	{0xF0, 2, 0xF4, registers.FlagsCleared, registers.FlagsCleared},
	{0xF0, 2, 0xF4, registers.FlagsFullSet, registers.FlagsFullSet},

	{0xF0, 4, 0xF0, registers.FlagsCleared, registers.FlagsCleared},
	{0xF0, 4, 0xF0, registers.FlagsFullSet, registers.FlagsFullSet},

	{0xF0, 6, 0xF0, registers.FlagsCleared, registers.FlagsCleared},
	{0xF0, 6, 0xF0, registers.FlagsFullSet, registers.FlagsFullSet},
}

func (s *unitTestSuite) TestSetBitInByte() {
	for _, c := range _setBitCases {
		name := fmt.Sprintf(
			"set bit %d with initialValue %#v, with flags (initial %#v, result %#v)",
			c.bit,
			c.initialValue,
			c.initialFlags,
			c.resultFlags,
		)

		s.Run(name, func() {
			s.Regs.F.Set(c.initialFlags)

			reg := registerslib.NewByte(c.initialValue)

			fn := s.factory.SetBitInByte(c.bit)(reg)
			result := fn()

			s.Equal(lib2.OpcodeSuccess(2, 8), result)
			s.Equal(c.resultValue, reg.Get())
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestSetBitInHLPtr_ValidCase() {
	for _, c := range _setBitCases {
		name := fmt.Sprintf(
			"set bit %d with initialValue %#v in hl ptr, with flags (initial %#v, result %#v)",
			c.bit,
			c.initialValue,
			c.initialFlags,
			c.resultFlags,
		)

		s.Run(name, func() {
			s.Regs.F.Set(c.initialFlags)

			s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(c.initialValue, nil)
			s.MockMMU.EXPECT().WriteByte(s.Regs.HL.Get(), c.resultValue).Return(nil)

			fn := s.factory.SetBitInHLPtr(c.bit)()
			result := fn()

			s.Equal(lib2.OpcodeSuccess(2, 16), result)
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestSetBitInHLPtr_InvalidHLRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.SetBitInHLPtr(0)()
	result := fn()

	s.Equal(lib2.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestSetBitInHLPtr_InvalidHLWrite() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0x00), nil)
	s.MockMMU.EXPECT().WriteByte(s.Regs.HL.Get(), uint8(0x01)).Return(testhelpers.ErrTesting1)

	fn := s.factory.SetBitInHLPtr(0)()
	result := fn()

	s.Equal(lib2.OpcodeError(testhelpers.ErrTesting1), result)
}
