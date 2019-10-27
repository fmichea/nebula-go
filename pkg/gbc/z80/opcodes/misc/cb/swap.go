package cb

import (
	"nebula-go/pkg/gbc/memory/registers"
)

func (f *Factory) innerSwap(value uint8) uint8 {
	newValue := ((value & 0x0F) << 4) | ((value & 0xF0) >> 4)
	return f.updateFlagsForCB(newValue, 0)
}

func (f *Factory) SwapByte(reg registers.Byte) cbopcode {
	return f.buildCBOpcodeByte(reg, f.innerSwap)
}

func (f *Factory) SwapHLPtr() cbopcode {
	return f.buildCBOpcodeHLPtr(f.innerSwap)
}
