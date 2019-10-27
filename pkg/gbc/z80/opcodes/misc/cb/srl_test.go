package cb

import (
	"fmt"
	"nebula-go/pkg/common/testhelpers"
	"nebula-go/pkg/gbc/memory/registers"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

var _srlCases = []struct {
	initialValue uint8
	resultValue  uint8

	initialFlags uint8
	resultFlags  uint8
}{
	{0x00, 0x00, z80lib.FlagsCleared, z80lib.ZF},
	{0x00, 0x00, z80lib.FlagsFullSet, z80lib.ZF},
	{0x0F, 0x07, z80lib.FlagsCleared, z80lib.CY},
	{0xF0, 0x78, z80lib.FlagsCleared, z80lib.FlagsCleared},
	{0xF0, 0x78, z80lib.FlagsFullSet, z80lib.FlagsCleared},
	{0xFF, 0x7F, z80lib.FlagsCleared, z80lib.CY},
}

func (s *unitTestSuite) TestSRLByte() {
	reg := registers.NewByte(0x00)
	fn := s.factory.SRLByte(reg)

	for _, c := range _srlCases {
		name := fmt.Sprintf(
			"srl in reg (initial = %#v, result = %#v) with flags (initial = %#v, result = %#v)",
			c.initialValue,
			c.resultValue,
			c.initialFlags,
			c.resultFlags,
		)

		s.Run(name, func() {
			reg.Set(c.initialValue)
			s.Regs.F.Set(c.initialFlags)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(2, 8), result)
			s.Equal(c.resultValue, reg.Get())
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestSRLHLPtr_ValidCase() {
	fn := s.factory.SRLHLPtr()

	for _, c := range _srlCases {
		name := fmt.Sprintf(
			"srl in hl ptr (initial = %#v, result = %#v) with flags (initial = %#v, result = %#v)",
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

			s.Equal(opcodeslib.OpcodeSuccess(2, 16), result)
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestSRLHLPtr_InvalidRead() {
	s.Regs.F.Set(z80lib.FlagsCleared)

	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.SRLHLPtr()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestSRLHLPtr_InvalidWrite() {
	s.Regs.F.Set(z80lib.FlagsCleared)

	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), nil)
	s.MockMMU.EXPECT().WriteByte(s.Regs.HL.Get(), uint8(0)).Return(testhelpers.ErrTesting1)

	fn := s.factory.SRLHLPtr()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
