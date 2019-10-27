package cb

import (
	"nebula-go/pkg/common/bitwise"
	"nebula-go/pkg/gbc/memory/registers"
)

func (f *Factory) sla(value uint8) uint8 {
	return f.updateFlagsForCB(value<<1, bitwise.HighBit8(value))
}

func (f *Factory) SLAByte(reg registers.Byte) cbopcode {
	return f.buildCBOpcodeByte(reg, f.sla)
}

func (f *Factory) SLAHLPtr() cbopcode {
	return f.buildCBOpcodeHLPtr(f.sla)
}
