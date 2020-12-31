package misc

import (
	"nebula-go/pkg/common/bitwise"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
)

func (f *Factory) RLCA() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		a := f.regs.A.Get()
		cy := bitwise.HighBit8(a)

		// NOTE: In the main z80 documentation, ZF is not affected, but based on test ROMs it seems like it is reset.
		f.regs.F.Set(registers.FlagsCleared)
		f.regs.F.CY.Set(cy)

		f.regs.A.Set((a << 1) | cy)

		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
