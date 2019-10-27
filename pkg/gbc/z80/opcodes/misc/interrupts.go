package misc

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) DI() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.regs.IME = false
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}

func (f *Factory) EI() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.regs.IME = true
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
