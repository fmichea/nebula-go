package misc

import (
	"nebula-go/pkg/common/bitwise"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) RRA() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		regs := f.regs

		a := regs.A.Get()
		cy := regs.F.CY.Get()

		// NOTE: see comment in RLA, there is some inconsistencies on the handling of flags here. Following the
		//  official z80 docs for now.
		regs.F.NE.SetBool(false)
		regs.F.HC.SetBool(false)
		regs.F.CY.Set(bitwise.LowBit8(a))

		regs.A.Set((a >> 1) | (cy << 7))

		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
