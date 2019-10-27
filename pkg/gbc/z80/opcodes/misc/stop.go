package misc

import opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"

func (f *Factory) Stop() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		// FIXME: implement this.
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
