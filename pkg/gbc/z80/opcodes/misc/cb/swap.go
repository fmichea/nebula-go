package cb

import registerslib "nebula-go/pkg/gbc/z80/registers/lib"

func (f *Factory) innerSwap(value uint8) uint8 {
	newValue := ((value & 0x0F) << 4) | ((value & 0xF0) >> 4)
	return f.updateFlagsForCB(newValue, 0)
}

func (f *Factory) SwapByte(reg registerslib.Byte) cbopcode {
	return f.buildCBOpcodeByte(reg, f.innerSwap)
}

func (f *Factory) SwapHLPtr() cbopcode {
	return f.buildCBOpcodeHLPtr(f.innerSwap)
}
