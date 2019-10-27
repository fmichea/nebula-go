package misc

import (
	"nebula-go/pkg/common/bitwise"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) RRCA() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		a := f.regs.A.Get()
		cy := bitwise.LowBit8(a)

		f.regs.F.NE.SetBool(false)
		f.regs.F.HC.SetBool(false)
		f.regs.F.CY.Set(cy)

		value := (a >> 1) | (cy << 7)
		f.regs.A.Set(value)

		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
