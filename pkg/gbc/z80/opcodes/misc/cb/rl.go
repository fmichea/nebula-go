package cb

import (
	"nebula-go/pkg/common/bitwise"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) rl(value uint8) uint8 {
	return f.updateFlagsForCB((value<<1)|f.regs.F.CY.Get(), bitwise.HighBit8(value))
}

func (f *Factory) RLByte(reg registerslib.Byte) cbopcode {
	return f.buildCBOpcodeByte(reg, f.rl)
}

func (f *Factory) RLHLPtr() cbopcode {
	return f.buildCBOpcodeHLPtr(f.rl)
}
