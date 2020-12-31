package misc

import (
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/memory/registers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (s *unitTestSuite) TestStop() {
	// FIXME: refactor a mockMMU to provide the registers mock here like in
	//  the graphics package without all of this work.
	var key1value uint8

	mmuRegs := &memory.Registers{
		KEY1: registers.NewKEY1Reg(&key1value),
	}

	stopOpcode := s.factory.Stop()

	s.Run("no switch request is noop, case 0", func() {
		s.MockMMU.EXPECT().Registers().Return(mmuRegs)

		key1value = 0x00

		result := stopOpcode()
		s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
		s.Equal(uint8(0x00), key1value)
	})

	s.Run("no switch request is noop, case 1", func() {
		s.MockMMU.EXPECT().Registers().Return(mmuRegs)

		key1value = 0x80

		result := stopOpcode()
		s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
		s.Equal(uint8(0x80), key1value)
	})

	s.Run("switch request is executed, case 0", func() {
		s.MockMMU.EXPECT().Registers().Return(mmuRegs)

		key1value = 0x01

		result := stopOpcode()
		s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
		s.Equal(uint8(0x80), key1value)
	})

	s.Run("switch request is executed, case 1", func() {
		s.MockMMU.EXPECT().Registers().Return(mmuRegs)

		key1value = 0x81

		result := stopOpcode()
		s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
		s.Equal(uint8(0x00), key1value)
	})
}
