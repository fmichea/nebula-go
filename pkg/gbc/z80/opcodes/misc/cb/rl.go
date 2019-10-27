package cb

import (
	"nebula-go/pkg/common/bitwise"
	"nebula-go/pkg/gbc/memory/registers"
)

func (f *Factory) rl(value uint8) uint8 {
	return f.updateFlagsForCB((value<<1)|f.regs.F.CY.Get(), bitwise.HighBit8(value))
}

func (f *Factory) RLByte(reg registers.Byte) cbopcode {
	return f.buildCBOpcodeByte(reg, f.rl)
}

func (f *Factory) RLHLPtr() cbopcode {
	return f.buildCBOpcodeHLPtr(f.rl)
}
