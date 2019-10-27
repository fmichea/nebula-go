package misc

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) CPL() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		regs := f.regs

		regs.A.Set(regs.A.Get() ^ 0xFF)

		regs.F.NE.SetBool(true)
		regs.F.HC.SetBool(true)

		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
