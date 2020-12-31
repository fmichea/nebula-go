package misc

import (
	"nebula-go/pkg/common/bitwise"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
)

func (f *Factory) RRA() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		regs := f.regs

		a := regs.A.Get()
		cy := regs.F.CY.Get()

		// NOTE: in the main z80 documentation, the ZF flag is not affected but based on the GameBoy opcode
		//  documentations, it is reset. Using the Gameboy behavior since it seems to be what ROMs use.
		regs.F.Set(registers.FlagsCleared)
		regs.F.CY.Set(bitwise.LowBit8(a))

		regs.A.Set((a >> 1) | (cy << 7))

		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
