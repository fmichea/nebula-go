package alu

import (
	"fmt"

	"nebula-go/pkg/common/testhelpers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (s *unitTestSuite) TestIncrementByte() {
	cases := []struct {
		initialValue uint8
		result       uint8
		initialFlags uint8
		resultFlags  uint8
	}{
		{0x00, 0x01, registers.FlagsCleared, registers.FlagsCleared},
		{0x00, 0x01, registers.FlagsFullSet, registers.CY},
		{0x0F, 0x10, registers.FlagsCleared, registers.HC},
		{0x0F, 0x10, registers.FlagsFullSet, registers.HC | registers.CY},
		{0x3F, 0x40, registers.FlagsCleared, registers.HC},
		{0x3F, 0x40, registers.FlagsFullSet, registers.HC | registers.CY},
		{0xFF, 0x00, registers.FlagsCleared, registers.ZF | registers.HC},
		{0xFF, 0x00, registers.FlagsFullSet, registers.ZF | registers.HC | registers.CY},
	}

	reg := registerslib.NewByte(0x00)
	fn := s.factory.IncrementByte(reg)

	for _, c := range cases {
		name := fmt.Sprintf("increment value %#v to %#v (initialFlags = %#v)", c.initialValue, c.result, c.initialFlags)

		s.Run(name, func() {
			reg.Set(c.initialValue)
			s.Regs.F.Set(c.initialFlags)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
			s.Equal(c.result, reg.Get())
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestIncrementDByte() {
	cases := []struct {
		initialValue uint16
		result       uint16
		flags        uint8
	}{
		{0x0000, 0x0001, registers.FlagsCleared},
		{0x0000, 0x0001, registers.FlagsFullSet},
		{0x000F, 0x0010, registers.FlagsCleared},
		{0x000F, 0x0010, registers.FlagsFullSet},
		{0x02FF, 0x0300, registers.FlagsCleared},
		{0x02FF, 0x0300, registers.FlagsFullSet},
		{0x9FFF, 0xA000, registers.FlagsCleared},
		{0x9FFF, 0xA000, registers.FlagsFullSet},
		{0xFFFF, 0x0000, registers.FlagsCleared},
		{0xFFFF, 0x0000, registers.FlagsFullSet},
	}

	reg := registerslib.NewDByte(0x00)
	fn := s.factory.IncrementDByte(reg)

	for _, c := range cases {
		name := fmt.Sprintf("increment value %#v to %#v (flags = %#v)", c.initialValue, c.result, c.flags)

		s.Run(name, func() {
			reg.Set(c.initialValue)
			s.Regs.F.Set(c.flags)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
			s.Equal(c.result, reg.Get())
			s.Equal(c.flags, s.Regs.F.Get())
		})
	}
}

func (s *unitTestSuite) TestIncrementHLPtr() {
	cases := []struct {
		initialValue uint8
		result       uint8
		initialFlags uint8
		resultFlags  uint8
	}{
		{0x00, 0x01, registers.FlagsCleared, registers.FlagsCleared},
		{0x00, 0x01, registers.FlagsFullSet, registers.CY},
		{0x0F, 0x10, registers.FlagsCleared, registers.HC},
		{0x0F, 0x10, registers.FlagsFullSet, registers.HC | registers.CY},
		{0x3F, 0x40, registers.FlagsCleared, registers.HC},
		{0x3F, 0x40, registers.FlagsFullSet, registers.HC | registers.CY},
		{0xFF, 0x00, registers.FlagsCleared, registers.ZF | registers.HC},
		{0xFF, 0x00, registers.FlagsFullSet, registers.ZF | registers.HC | registers.CY},
	}

	fn := s.factory.IncrementHLPtr()

	for _, c := range cases {
		name := fmt.Sprintf("increment value %#v to %#v (initialFlags = %#v)", c.initialValue, c.result, c.initialFlags)

		s.Run(name, func() {
			s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(c.initialValue, nil)
			s.MockMMU.EXPECT().WriteByte(s.Regs.HL.Get(), c.result).Return(nil)

			s.Regs.F.Set(c.initialFlags)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(1, 12), result)
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestIncrementHLPtr_InvalidHLRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.IncrementHLPtr()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestIncrementHLPtr_InvalidHLWrite() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), nil)
	s.MockMMU.EXPECT().WriteByte(s.Regs.HL.Get(), uint8(1)).Return(testhelpers.ErrTesting1)

	fn := s.factory.IncrementHLPtr()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
