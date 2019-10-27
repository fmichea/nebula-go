package cb

import (
	"nebula-go/pkg/common/bitwise"
	"nebula-go/pkg/gbc/memory/registers"
)

func (f *Factory) sra(value uint8) uint8 {
	newValue := (value & 0x80) | (value >> 1)
	return f.updateFlagsForCB(newValue, bitwise.LowBit8(value))
}

func (f *Factory) SRAByte(reg registers.Byte) cbopcode {
	return f.buildCBOpcodeByte(reg, f.sra)
}

func (f *Factory) SRAHLPtr() cbopcode {
	return f.buildCBOpcodeHLPtr(f.sra)
}
