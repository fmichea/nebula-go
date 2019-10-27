package misc

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) Noop() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
