package cb

import (
	"nebula-go/pkg/common/bitwise"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) sla(value uint8) uint8 {
	return f.updateFlagsForCB(value<<1, bitwise.HighBit8(value))
}

func (f *Factory) SLAByte(reg registerslib.Byte) cbopcode {
	return f.buildCBOpcodeByte(reg, f.sla)
}

func (f *Factory) SLAHLPtr() cbopcode {
	return f.buildCBOpcodeHLPtr(f.sla)
}
