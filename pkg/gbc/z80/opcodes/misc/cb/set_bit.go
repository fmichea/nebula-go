package cb

import registerslib "nebula-go/pkg/gbc/z80/registers/lib"

func (f *Factory) setBit(bit uint8) func(uint8) uint8 {
	return func(value uint8) uint8 {
		return value | (0x1 << bit)
	}
}

func (f *Factory) SetBitInByte(bit uint8) cbbytefunc {
	return func(reg registerslib.Byte) cbopcode {
		return f.buildCBOpcodeByte(reg, f.setBit(bit))
	}
}

func (f *Factory) SetBitInHLPtr(bit uint8) cbhlptrfunc {
	return func() cbopcode {
		return f.buildCBOpcodeHLPtr(f.setBit(bit))
	}
}
