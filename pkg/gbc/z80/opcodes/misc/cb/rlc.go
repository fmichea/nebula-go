package cb

import (
	"nebula-go/pkg/common/bitwise"
	"nebula-go/pkg/gbc/memory/registers"
)

func (f *Factory) rlc(value uint8) uint8 {
	cy := bitwise.HighBit8(value)
	return f.updateFlagsForCB((value<<1)|cy, cy)
}

func (f *Factory) RLCByte(reg registers.Byte) cbopcode {
	return f.buildCBOpcodeByte(reg, f.rlc)
}

func (f *Factory) RLCHLPtr() cbopcode {
	return f.buildCBOpcodeHLPtr(f.rlc)
}
