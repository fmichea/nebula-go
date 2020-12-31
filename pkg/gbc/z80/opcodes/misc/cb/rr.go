package cb

import (
	"nebula-go/pkg/common/bitwise"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) rr(value uint8) uint8 {
	newValue := (f.regs.F.CY.Get() << 7) | (value >> 1)
	return f.updateFlagsForCB(newValue, bitwise.LowBit8(value))
}

func (f *Factory) RRByte(reg registerslib.Byte) cbopcode {
	return f.buildCBOpcodeByte(reg, f.rr)
}

func (f *Factory) RRHLPtr() cbopcode {
	return f.buildCBOpcodeHLPtr(f.rr)
}
