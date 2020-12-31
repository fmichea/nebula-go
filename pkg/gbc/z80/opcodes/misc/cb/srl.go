package cb

import (
	"nebula-go/pkg/common/bitwise"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) srl(value uint8) uint8 {
	return f.updateFlagsForCB(value>>1, bitwise.LowBit8(value))
}

func (f *Factory) SRLByte(reg registerslib.Byte) cbopcode {
	return f.buildCBOpcodeByte(reg, f.srl)
}

func (f *Factory) SRLHLPtr() cbopcode {
	return f.buildCBOpcodeHLPtr(f.srl)
}
