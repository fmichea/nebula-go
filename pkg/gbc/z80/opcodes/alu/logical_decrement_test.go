package alu

import (
	"fmt"
	"nebula-go/pkg/common/testhelpers"
	"nebula-go/pkg/gbc/memory/registers"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) TestDecrementByte() {
	cases := []struct {
		initialValue uint8
		result       uint8
		initialFlags uint8
		resultFlags  uint8
	}{
		{0x01, 0x00, z80lib.FlagsCleared, z80lib.ZF},
		{0x01, 0x00, z80lib.FlagsFullSet, z80lib.ZF | z80lib.CY},
		{0x02, 0x01, z80lib.FlagsCleared, z80lib.FlagsCleared},
		{0x02, 0x01, z80lib.FlagsFullSet, z80lib.CY},
		{0x10, 0x0F, z80lib.FlagsCleared, z80lib.HC},
		{0x10, 0x0F, z80lib.FlagsFullSet, z80lib.HC | z80lib.CY},
		{0x40, 0x3F, z80lib.FlagsCleared, z80lib.HC},
		{0x40, 0x3F, z80lib.FlagsFullSet, z80lib.HC | z80lib.CY},
		{0x00, 0xFF, z80lib.FlagsCleared, z80lib.HC},
		{0x00, 0xFF, z80lib.FlagsFullSet, z80lib.HC | z80lib.CY},
	}

	reg := registers.NewByte(0x00)
	fn := s.factory.DecrementByte(reg)

	for _, c := range cases {
		name := fmt.Sprintf("decrement value %#v to %#v (initialFlags = %#v)", c.initialValue, c.result, c.initialFlags)

		s.Run(name, func() {
			reg.Set(c.initialValue)
			s.Regs.F.Set(c.initialFlags)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
			s.Equal(c.result, reg.Get())
			s.EqualFlags(c.resultFlags | z80lib.NE)
		})
	}
}

func (s *unitTestSuite) TestDecrementDByte() {
	cases := []struct {
		initialValue uint16
		result       uint16
		flags        uint8
	}{
		{0x0002, 0x0001, z80lib.FlagsCleared},
		{0x0002, 0x0001, z80lib.FlagsFullSet},
		{0x0010, 0x000F, z80lib.FlagsCleared},
		{0x0010, 0x000F, z80lib.FlagsFullSet},
		{0x0300, 0x02FF, z80lib.FlagsCleared},
		{0x0300, 0x02FF, z80lib.FlagsFullSet},
		{0xA000, 0x9FFF, z80lib.FlagsCleared},
		{0xA000, 0x9FFF, z80lib.FlagsFullSet},
		{0x0000, 0xFFFF, z80lib.FlagsCleared},
		{0x0000, 0xFFFF, z80lib.FlagsFullSet},
	}

	reg := registers.NewDByte(0x0000)
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
		{0x01, 0x00, z80lib.FlagsCleared, z80lib.ZF},
		{0x01, 0x00, z80lib.FlagsFullSet, z80lib.ZF | z80lib.CY},
		{0x02, 0x01, z80lib.FlagsCleared, z80lib.FlagsCleared},
		{0x02, 0x01, z80lib.FlagsFullSet, z80lib.CY},
		{0x10, 0x0F, z80lib.FlagsCleared, z80lib.HC},
		{0x10, 0x0F, z80lib.FlagsFullSet, z80lib.HC | z80lib.CY},
		{0x40, 0x3F, z80lib.FlagsCleared, z80lib.HC},
		{0x40, 0x3F, z80lib.FlagsFullSet, z80lib.HC | z80lib.CY},
		{0x00, 0xFF, z80lib.FlagsCleared, z80lib.HC},
		{0x00, 0xFF, z80lib.FlagsFullSet, z80lib.HC | z80lib.CY},
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
			s.EqualFlags(c.resultFlags | z80lib.NE)
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
