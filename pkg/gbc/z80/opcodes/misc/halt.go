package misc

import opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"

// FIXME: See if there isn't something else to do here?
func (f *Factory) Halt() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.regs.HaltMode = true
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
