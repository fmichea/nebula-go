package cb

import (
	"fmt"
	"nebula-go/pkg/common/testhelpers"
	"nebula-go/pkg/gbc/memory/registers"
	z80_lib "nebula-go/pkg/gbc/z80/lib"
	lib2 "nebula-go/pkg/gbc/z80/opcodes/lib"
)

var _testBitCases = []cbBitOpTestCase{
	{0xF0, 0, 0xF0, z80_lib.FlagsCleared, z80_lib.ZF | z80_lib.HC},
	{0xF0, 0, 0xF0, z80_lib.FlagsFullSet, z80_lib.ZF | z80_lib.HC | z80_lib.CY},

	{0xF0, 2, 0xF0, z80_lib.FlagsCleared, z80_lib.ZF | z80_lib.HC},
	{0xF0, 2, 0xF0, z80_lib.FlagsFullSet, z80_lib.ZF | z80_lib.HC | z80_lib.CY},

	{0xF0, 4, 0xF0, z80_lib.FlagsCleared, z80_lib.HC},
	{0xF0, 4, 0xF0, z80_lib.FlagsFullSet, z80_lib.HC | z80_lib.CY},

	{0xF0, 6, 0xF0, z80_lib.FlagsCleared, z80_lib.HC},
	{0xF0, 6, 0xF0, z80_lib.FlagsFullSet, z80_lib.HC | z80_lib.CY},
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

			reg := registers.NewByte(c.initialValue)

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

			s.Equal(lib2.OpcodeSuccess(2, 16), result)
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
