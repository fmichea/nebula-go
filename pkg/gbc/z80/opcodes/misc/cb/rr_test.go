package cb

import (
	"fmt"

	"nebula-go/pkg/common/testhelpers"
	lib2 "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

var _rrCases = []struct {
	initialValue uint8
	resultValue  uint8

	initialFlags uint8
	resultFlags  uint8
}{
	{0x00, 0x00, registers.FlagsCleared, registers.ZF},
	{0x00, 0x80, registers.FlagsFullSet, registers.FlagsCleared},
	{0x00, 0x00, registers.ZF | registers.NE | registers.HC, registers.ZF},
	{0xF0, 0x78, registers.FlagsCleared, registers.FlagsCleared},
	{0xF1, 0x78, registers.FlagsCleared, registers.CY},
	{0xF1, 0xF8, registers.CY, registers.CY},
}

func (s *unitTestSuite) TestRRByte() {
	reg := registerslib.NewByte(0x00)
	fn := s.factory.RRByte(reg)

	for _, c := range _rrCases {
		name := fmt.Sprintf(
			"rr in reg (initial = %#v, result = %#v) with flags (initial = %#v, result = %#v)",
			c.initialValue,
			c.resultValue,
			c.initialFlags,
			c.resultFlags,
		)

		s.Run(name, func() {
			reg.Set(c.initialValue)
			s.Regs.F.Set(c.initialFlags)

			result := fn()

			s.Equal(lib2.OpcodeSuccess(2, 8), result)
			s.Equal(c.resultValue, reg.Get())
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestRRHLPtr_ValidCase() {
	fn := s.factory.RRHLPtr()

	for _, c := range _rrCases {
		name := fmt.Sprintf(
			"rr in hl ptr (initial = %#v, result = %#v) with flags (initial = %#v, result = %#v)",
			c.initialValue,
			c.resultValue,
			c.initialFlags,
			c.resultFlags,
		)

		s.Run(name, func() {
			s.Regs.F.Set(c.initialFlags)

			s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(c.initialValue, nil)
			s.MockMMU.EXPECT().WriteByte(s.Regs.HL.Get(), c.resultValue).Return(nil)

			result := fn()

			s.Equal(lib2.OpcodeSuccess(2, 16), result)
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestRRHLPtr_InvalidRead() {
	s.Regs.F.Set(registers.FlagsCleared)

	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.RRHLPtr()
	result := fn()

	s.Equal(lib2.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestRRHLPtr_InvalidWrite() {
	s.Regs.F.Set(registers.FlagsCleared)

	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), nil)
	s.MockMMU.EXPECT().WriteByte(s.Regs.HL.Get(), uint8(0)).Return(testhelpers.ErrTesting1)

	fn := s.factory.RRHLPtr()
	result := fn()

	s.Equal(lib2.OpcodeError(testhelpers.ErrTesting1), result)
}
