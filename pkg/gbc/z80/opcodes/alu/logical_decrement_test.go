package alu

import (
	"fmt"

	"nebula-go/pkg/common/testhelpers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (s *unitTestSuite) TestDecrementByte() {
	cases := []struct {
		initialValue uint8
		result       uint8
		initialFlags uint8
		resultFlags  uint8
	}{
		{0x01, 0x00, registers.FlagsCleared, registers.ZF},
		{0x01, 0x00, registers.FlagsFullSet, registers.ZF | registers.CY},
		{0x02, 0x01, registers.FlagsCleared, registers.FlagsCleared},
		{0x02, 0x01, registers.FlagsFullSet, registers.CY},
		{0x10, 0x0F, registers.FlagsCleared, registers.HC},
		{0x10, 0x0F, registers.FlagsFullSet, registers.HC | registers.CY},
		{0x40, 0x3F, registers.FlagsCleared, registers.HC},
		{0x40, 0x3F, registers.FlagsFullSet, registers.HC | registers.CY},
		{0x00, 0xFF, registers.FlagsCleared, registers.HC},
		{0x00, 0xFF, registers.FlagsFullSet, registers.HC | registers.CY},
	}

	reg := registerslib.NewByte(0x00)
	fn := s.factory.DecrementByte(reg)

	for _, c := range cases {
		name := fmt.Sprintf("decrement value %#v to %#v (initialFlags = %#v)", c.initialValue, c.result, c.initialFlags)

		s.Run(name, func() {
			reg.Set(c.initialValue)
			s.Regs.F.Set(c.initialFlags)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
			s.Equal(c.result, reg.Get())
			s.EqualFlags(c.resultFlags | registers.NE)
		})
	}
}

func (s *unitTestSuite) TestDecrementDByte() {
	cases := []struct {
		initialValue uint16
		result       uint16
		flags        uint8
	}{
		{0x0002, 0x0001, registers.FlagsCleared},
		{0x0002, 0x0001, registers.FlagsFullSet},
		{0x0010, 0x000F, registers.FlagsCleared},
		{0x0010, 0x000F, registers.FlagsFullSet},
		{0x0300, 0x02FF, registers.FlagsCleared},
		{0x0300, 0x02FF, registers.FlagsFullSet},
		{0xA000, 0x9FFF, registers.FlagsCleared},
		{0xA000, 0x9FFF, registers.FlagsFullSet},
		{0x0000, 0xFFFF, registers.FlagsCleared},
		{0x0000, 0xFFFF, registers.FlagsFullSet},
	}

	reg := registerslib.NewDByte(0x0000)
	fn := s.factory.DecrementDByte(reg)

	for _, c := range cases {
		name := fmt.Sprintf("decrement value %#v to %#v (flags = %#v)", c.initialValue, c.result, c.flags)

		s.Run(name, func() {
			reg.Set(c.initialValue)
			s.Regs.F.Set(c.flags)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
			s.Equal(c.result, reg.Get())
			s.EqualFlags(c.flags)
		})
	}
}

func (s *unitTestSuite) TestDecrementHLPtr() {
	cases := []struct {
		initialValue uint8
		result       uint8
		initialFlags uint8
		resultFlags  uint8
	}{
		{0x01, 0x00, registers.FlagsCleared, registers.ZF},
		{0x01, 0x00, registers.FlagsFullSet, registers.ZF | registers.CY},
		{0x02, 0x01, registers.FlagsCleared, registers.FlagsCleared},
		{0x02, 0x01, registers.FlagsFullSet, registers.CY},
		{0x10, 0x0F, registers.FlagsCleared, registers.HC},
		{0x10, 0x0F, registers.FlagsFullSet, registers.HC | registers.CY},
		{0x40, 0x3F, registers.FlagsCleared, registers.HC},
		{0x40, 0x3F, registers.FlagsFullSet, registers.HC | registers.CY},
		{0x00, 0xFF, registers.FlagsCleared, registers.HC},
		{0x00, 0xFF, registers.FlagsFullSet, registers.HC | registers.CY},
	}

	fn := s.factory.DecrementHLPtr()

	for _, c := range cases {
		name := fmt.Sprintf("decrement value %#v to %#v (initialFlags = %#v)", c.initialValue, c.result, c.initialFlags)

		s.Run(name, func() {
			s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(c.initialValue, nil)
			s.MockMMU.EXPECT().WriteByte(s.Regs.HL.Get(), c.result).Return(nil)

			s.Regs.F.Set(c.initialFlags)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(1, 12), result)
			s.EqualFlags(c.resultFlags | registers.NE)
		})
	}
}

func (s *unitTestSuite) TestDecrementHLPtr_InvalidHLRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(0), testhelpers.ErrTesting1)

	fn := s.factory.DecrementHLPtr()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestDecrementHLPtr_InvalidHLWrite() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.HL.Get()).Return(uint8(2), nil)
	s.MockMMU.EXPECT().WriteByte(s.Regs.HL.Get(), uint8(1)).Return(testhelpers.ErrTesting1)

	fn := s.factory.DecrementHLPtr()
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}
