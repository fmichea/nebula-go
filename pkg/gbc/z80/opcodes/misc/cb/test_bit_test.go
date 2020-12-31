package cb

import (
	"fmt"

	"nebula-go/pkg/common/testhelpers"
	lib2 "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

var _testBitCases = []cbBitOpTestCase{
	{0xF0, 0, 0xF0, registers.FlagsCleared, registers.ZF | registers.HC},
	{0xF0, 0, 0xF0, registers.FlagsFullSet, registers.ZF | registers.HC | registers.CY},

	{0xF0, 2, 0xF0, registers.FlagsCleared, registers.ZF | registers.HC},
	{0xF0, 2, 0xF0, registers.FlagsFullSet, registers.ZF | registers.HC | registers.CY},

	{0xF0, 4, 0xF0, registers.FlagsCleared, registers.HC},
	{0xF0, 4, 0xF0, registers.FlagsFullSet, registers.HC | registers.CY},

	{0xF0, 6, 0xF0, registers.FlagsCleared, registers.HC},
	{0xF0, 6, 0xF0, registers.FlagsFullSet, registers.HC | registers.CY},
}

func (s *unitTestSuite) TestTestBitInByte() {
	for _, c := range _testBitCases {
		name := fmt.Sprintf(
			"test bit %d with initialValue %#v, with flags (initial %#v, result %#v)",
			c.bit,
			c.initialValue,
			c.initialFlags,
			c.resultFlags,
		)

		s.Run(name, func() {
			s.Regs.F.Set(c.initialFlags)

			reg := registerslib.NewByte(c.initialValue)

			fn := s.factory.TestBitInByte(c.bit)(reg)
			result := fn()

			s.Equal(lib2.OpcodeSuccess(2, 8), result)
			s.Equal(c.resultValue, reg.Get())
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestTestBitInHLPtr_ValidCase() {
	for _, c := range _testBitCases {
		name := fmt.Sprintf(
			"test bit %d with initialValue %#v in hl ptr, with flags (initial %#v, result %#v)",
			c.bit,
			c.initialValue,
			c.initialFlags,
			c.resultFlags,
		)

		s.Run(name, func() {
			s.Regs.F.Set(c.initialFlags)

			s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(c.initialValue, nil)

			fn := s.factory.TestBitInHLPtr(c.bit)()
			result := fn()

			s.Equal(lib2.OpcodeSuccess(2, 12), result)
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestTestBitInHLPtr_InvalidHLRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.TestBitInHLPtr(0)()
	result := fn()

	s.Equal(lib2.OpcodeError(testhelpers.ErrTesting1), result)
}
