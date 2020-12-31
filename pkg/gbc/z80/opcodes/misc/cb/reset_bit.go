package cb

import (
	"nebula-go/pkg/common/bitwise"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) resetBit(bit uint8) func(uint8) uint8 {
	return func(value uint8) uint8 {
		return bitwise.InverseMask8(value, 0x1<<bit)
	}
}

func (f *Factory) ResetBitInByte(bit uint8) cbbytefunc {
	return func(reg registerslib.Byte) cbopcode {
		return f.buildCBOpcodeByte(reg, f.resetBit(bit))
	}
}

func (f *Factory) ResetBitInHLPtr(bit uint8) cbhlptrfunc {
	return func() cbopcode {
		return f.buildCBOpcodeHLPtr(f.resetBit(bit))
	}
}
