package misc

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) CCF() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		cy := f.regs.F.CY.Get()

		f.regs.F.NE.SetBool(false)
		f.regs.F.HC.Set(cy)
		f.regs.F.CY.SetBool(cy == 0)
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
