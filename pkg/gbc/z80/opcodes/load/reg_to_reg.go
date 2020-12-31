package load

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) ByteToByte(dest, source registerslib.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		dest.Set(source.Get())
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
