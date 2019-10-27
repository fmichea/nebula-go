package cb

import (
	"nebula-go/pkg/common/bitwise"
	"nebula-go/pkg/gbc/memory/registers"
)

func (f *Factory) rrc(value uint8) uint8 {
	cy := bitwise.LowBit8(value)
	newValue := (cy << 7) | (value >> 1)
	return f.updateFlagsForCB(newValue, cy)
}

func (f *Factory) RRCByte(reg registers.Byte) cbopcode {
	return f.buildCBOpcodeByte(reg, f.rrc)
}

func (f *Factory) RRCHLPtr() cbopcode {
	return f.buildCBOpcodeHLPtr(f.rrc)
}
