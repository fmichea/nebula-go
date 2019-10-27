package cb

import (
	"nebula-go/pkg/common/bitwise"
	"nebula-go/pkg/gbc/memory/registers"
)

func (f *Factory) srl(value uint8) uint8 {
	return f.updateFlagsForCB(value>>1, bitwise.LowBit8(value))
}

func (f *Factory) SRLByte(reg registers.Byte) cbopcode {
	return f.buildCBOpcodeByte(reg, f.srl)
}

func (f *Factory) SRLHLPtr() cbopcode {
	return f.buildCBOpcodeHLPtr(f.srl)
}
