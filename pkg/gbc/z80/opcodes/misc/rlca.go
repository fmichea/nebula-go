package misc

import (
	"nebula-go/pkg/common/bitwise"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) RLCA() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		a := f.regs.A.Get()
		cy := bitwise.HighBit8(a)

		// NOTE: inconsistency found here again (see comment in rla), following the main z80 documentation.
		f.regs.F.NE.SetBool(false)
		f.regs.F.HC.SetBool(false)
		f.regs.F.CY.Set(cy)

		f.regs.A.Set((a << 1) | cy)

		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
