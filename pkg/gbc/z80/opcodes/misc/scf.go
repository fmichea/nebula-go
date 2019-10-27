package misc

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) SCF() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.regs.F.HC.SetBool(false)
		f.regs.F.NE.SetBool(false)
		f.regs.F.CY.SetBool(true)
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
