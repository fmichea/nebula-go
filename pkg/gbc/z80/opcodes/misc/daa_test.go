package misc

import (
	"fmt"

	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
)

func (s *unitTestSuite) TestDAA() {
	cases := []struct {
		initialValue uint8
		resultValue  uint8

		initialFlags uint8
		resultFlags  uint8
	}{
		// Additions:
		// 0b0000 (0x00) + 0b0001 (0x01) = 0b0001 (0x01)
		{0x01, 0x01, registers.FlagsCleared, registers.FlagsCleared},
		// 0b0101 (0x05) + 0b0101 (0x05) = 0b1010 (0x0A), should be 0b00010000 (0x10)
		{0x0A, 0x10, registers.FlagsCleared, registers.FlagsCleared},
		// 0b1001 (0x09) + 0b1001 (0x09) = 0b00010010 (0x12), should be 0b00011000 (0x18)
		{0x12, 0x18, registers.HC, registers.FlagsCleared},
		// 0b10011001 (0x99) + 0b00000001 (0x01) = 0b10011010 (0x9A), should be 0b00000000 (0x00, really 0x100)
		{0x9A, 0x00, registers.FlagsCleared, registers.ZF | registers.CY},
		// 0b10011001 (0x99) + 0b10011001 (0x99) = 0b00110010 (0x32), should be 0b10011000 (0x98, really 0x198)
		{0x32, 0x98, registers.HC | registers.CY, registers.CY},

		// Subtractions:
		// 0b0010 (0x02) - 0b0001 (0x01) = 0b0001 (0x01)
		{0x01, 0x01, registers.NE, registers.NE},
		// 0b01000111 (0x47) - 0b00101000 (0x28) = 0b00011111 (0x1F), should be 0b00011001 (0x19)
		{0x1F, 0x19, registers.NE | registers.HC, registers.NE},
		// 0b00000000 (0x00, really 0x100) - 0b1 (0x01) = 0b11111111 (0xFF), should be 0b10011001 (0x99)
		{0xFF, 0x99, registers.NE | registers.HC | registers.CY, registers.NE | registers.CY},
	}

	fn := s.factory.DAA()

	for _, c := range cases {
		name := fmt.Sprintf(
			"daa a reg (initial = %#v, result = %#v) and flags (initial = %#v, result = %#v)",
			c.initialValue,
			c.resultValue,
			c.initialFlags,
			c.resultFlags,
		)

		s.Run(name, func() {
			s.Regs.A.Set(c.initialValue)
			s.Regs.F.Set(c.initialFlags)

			result := fn()
			s.Equal(opcodeslib.OpcodeSuccess(1, 4), result)
			s.Equal(c.resultValue, s.Regs.A.Get())
			s.EqualFlags(c.resultFlags)
		})
	}
}
