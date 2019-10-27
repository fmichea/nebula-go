package cb

import (
	"nebula-go/pkg/common/bitwise"
	"nebula-go/pkg/gbc/memory/registers"
)

func (f *Factory) rr(value uint8) uint8 {
	newValue := (f.regs.F.CY.Get() << 7) | (value >> 1)
	return f.updateFlagsForCB(newValue, bitwise.LowBit8(value))
}

func (f *Factory) RRByte(reg registers.Byte) cbopcode {
	return f.buildCBOpcodeByte(reg, f.rr)
}

func (f *Factory) RRHLPtr() cbopcode {
	return f.buildCBOpcodeHLPtr(f.rr)
}
