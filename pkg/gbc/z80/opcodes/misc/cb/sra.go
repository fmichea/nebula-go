package cb

import (
	"nebula-go/pkg/common/bitwise"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) sra(value uint8) uint8 {
	newValue := (value & 0x80) | (value >> 1)
	return f.updateFlagsForCB(newValue, bitwise.LowBit8(value))
}

func (f *Factory) SRAByte(reg registerslib.Byte) cbopcode {
	return f.buildCBOpcodeByte(reg, f.sra)
}

func (f *Factory) SRAHLPtr() cbopcode {
	return f.buildCBOpcodeHLPtr(f.sra)
}
