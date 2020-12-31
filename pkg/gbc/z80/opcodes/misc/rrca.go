package misc

import (
	"nebula-go/pkg/common/bitwise"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
)

func (f *Factory) RRCA() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		a := f.regs.A.Get()
		cy := bitwise.LowBit8(a)

		// NOTE: in the main z80 documentation, the ZF flag is not affected but based on the GameBoy opcode
		//documentations, it is reset. Using the Gameboy behavior since it seems to be what ROMs use.
		f.regs.F.Set(registers.FlagsCleared)
		f.regs.F.CY.Set(cy)

		value := (a >> 1) | (cy << 7)
		f.regs.A.Set(value)

		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
