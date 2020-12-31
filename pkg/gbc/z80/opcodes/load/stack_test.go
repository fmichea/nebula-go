package load

import (
	"fmt"

	"nebula-go/pkg/common/testhelpers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (s *unitTestSuite) TestPushDByte_ValidCase() {
	value := uint16(0xABCD)
	adjustedSP := s.Regs.SP.Get() - 2

	reg := registerslib.NewDByte(value)

	s.MockMMU.EXPECT().WriteDByte(adjustedSP, value).Return(nil)

	fn := s.factory.PushDByte(reg)
	result := fn()

	s.Equal(opcodeslib.OpcodeSuccess(1, 16), result)
	s.Equal(value, reg.Get())
	s.Equal(adjustedSP, s.Regs.SP.Get())
}

func (s *unitTestSuite) TestPushDByte_InvalidWrite() {
	value := uint16(0xABCD)
	adjustedSP := s.Regs.SP.Get() - 2

	reg := registerslib.NewDByte(value)

	s.MockMMU.EXPECT().WriteDByte(adjustedSP, value).Return(testhelpers.ErrTesting1)

	fn := s.factory.PushDByte(reg)
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestPopDByte_ValidCase() {
	value := uint16(0xABCD)
	sp := s.Regs.SP.Get()

	reg := registerslib.NewDByte(0x0000)

	s.MockMMU.EXPECT().ReadDByte(sp).Return(value, nil)

	fn := s.factory.PopDByte(reg)
	result := fn()

	s.Equal(opcodeslib.OpcodeSuccess(1, 12), result)
	s.Equal(value, reg.Get())
	s.Equal(sp+2, s.Regs.SP.Get())
}

func (s *unitTestSuite) TestPopDByte_InvalidRead() {
	sp := s.Regs.SP.Get()

	reg := registerslib.NewDByte(0x0000)

	s.MockMMU.EXPECT().ReadDByte(sp).Return(uint16(0), testhelpers.ErrTesting1)

	fn := s.factory.PopDByte(reg)
	result := fn()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestLoadSPToAddress_ValidCase() {
	addr := uint16(0xABCD)

	pc := s.Regs.PC
	sp := s.Regs.SP.Get()

	s.MockMMU.EXPECT().ReadDByte(pc+1).Return(addr, nil)
	s.MockMMU.EXPECT().WriteDByte(addr, sp).Return(nil)

	result := s.factory.SPToAddress()()

	s.Equal(opcodeslib.OpcodeSuccess(3, 20), result)
}

func (s *unitTestSuite) TestLoadSPToAddress_InvalidRead() {
	addr := uint16(0xABCD)

	pc := s.Regs.PC

	s.MockMMU.EXPECT().ReadDByte(pc+1).Return(addr, testhelpers.ErrTesting1)

	result := s.factory.SPToAddress()()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestLoadSPToAddress_InvalidWrite() {
	addr := uint16(0xABCD)

	pc := s.Regs.PC
	sp := s.Regs.SP.Get()

	s.MockMMU.EXPECT().ReadDByte(pc+1).Return(addr, nil)
	s.MockMMU.EXPECT().WriteDByte(addr, sp).Return(testhelpers.ErrTesting1)

	result := s.factory.SPToAddress()()

	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestSPR8ToHL_ValidCase() {
	cases := []struct {
		initialSP uint16
		r8        uint8
		resultHL  uint16

		initialFlags uint8
		resultFlags  uint8
	}{
		{0xABCD, 0x00, 0xABCD, registers.FlagsFullSet, registers.FlagsCleared},

		{0xABCD, 0x01, 0xABCE, registers.FlagsCleared, registers.FlagsCleared},
		{0xABCD, 0x01, 0xABCE, registers.FlagsFullSet, registers.FlagsCleared},
		{0xABCD, 0x04, 0xABD1, registers.FlagsCleared, registers.HC},
		{0xABFD, 0x10, 0xAC0D, registers.FlagsCleared, registers.CY},
		{0xABCD, 0x73, 0xAC40, registers.FlagsCleared, registers.HC | registers.CY},

		// Addition with negative:
		{0xAB1D, 0x85, 0xAAA2, registers.FlagsCleared, registers.HC},
		{0xABF0, 0xCE, 0xABBE, registers.FlagsCleared, registers.CY},
		{0xABF5, 0xCE, 0xABC3, registers.FlagsCleared, registers.HC | registers.CY},
	}

	fn := s.factory.SPR8ToHL()

	for _, c := range cases {
		name := fmt.Sprintf(
			"load sp+r8 to hl: sp = %#v, hl = %#v and r8 = %#v",
			c.initialSP,
			c.resultHL,
			c.r8,
		)

		s.Run(name, func() {
			s.Regs.F.Set(c.initialFlags)
			s.Regs.SP.Set(c.initialSP)

			s.MockMMU.EXPECT().ReadByte(s.Regs.PC+1).Return(c.r8, nil)

			result := fn()

			s.Equal(opcodeslib.OpcodeSuccess(2, 12), result)
			s.Equal(c.initialSP, s.Regs.SP.Get())
			s.Equal(c.resultHL, s.Regs.HL.Get())
			s.EqualFlags(c.resultFlags)
		})
	}
}

func (s *unitTestSuite) TestSPR8ToHL_InvalidRead() {
	s.MockMMU.EXPECT().ReadByte(s.Regs.PC+1).Return(uint8(0), testhelpers.ErrTesting1)

	result := s.factory.SPR8ToHL()()
	s.Equal(opcodeslib.OpcodeError(testhelpers.ErrTesting1), result)
}

func (s *unitTestSuite) TestHLToSP() {
	v16 := uint16(0xABCD)

	s.Regs.SP.Set(0x0000)
	s.Regs.HL.Set(v16)

	result := s.factory.HLToSP()()
	s.Equal(opcodeslib.OpcodeSuccess(1, 8), result)
	s.Equal(v16, s.Regs.SP.Get())
}
