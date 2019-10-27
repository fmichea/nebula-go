package load

import (
	"nebula-go/pkg/gbc/memory/registers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) ByteToByte(dest, source registers.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		dest.Set(source.Get())
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
