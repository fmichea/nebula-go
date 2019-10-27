package cb

import (
	"fmt"
	"nebula-go/pkg/common/testhelpers"
	"nebula-go/pkg/gbc/memory/registers"
	z80_lib "nebula-go/pkg/gbc/z80/lib"
	lib2 "nebula-go/pkg/gbc/z80/opcodes/lib"
)

var _rrCases = []struct {
	initialValue uint8
	resultValue  uint8

	initialFlags uint8
	resultFlags  uint8
}{
	{0x00, 0x00, z80_lib.FlagsCleared, z80_lib.ZF},
	{0x00, 0x80, z80_lib.FlagsFullSet, z80_lib.FlagsCleared},
	{0x00, 0x00, z80_lib.ZF | z80_lib.NE | z80_lib.HC, z80_lib.ZF},
	{0xF0, 0x78, z80_lib.FlagsCleared, z80_lib.FlagsCleared},
	{0xF1, 0x78, z80_lib.FlagsCleared, z80_lib.CY},
	{0xF1, 0xF8, z80_lib.CY, z80_lib.CY},
}

func (s *unitTestSuite) TestRRByte() {
	reg := registers.NewByte(0x00)
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
	s.Regs.F.Set(z80_lib.FlagsCleared)

	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.RRHLPtr()
	result := fn()

	s.Equal(lib2.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestRRHLPtr_InvalidWrite() {
	s.Regs.F.Set(z80_lib.FlagsCleared)

	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), nil)
	s.MockMMU.EXPECT().WriteByte(s.Regs.HL.Get(), uint8(0)).Return(testhelpers.ErrTesting1)

	fn := s.factory.RRHLPtr()
	result := fn()

	s.Equal(lib2.OpcodeError(testhelpers.ErrTesting1), result)
}
