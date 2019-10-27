package controlflow

import (
	"nebula-go/pkg/gbc/memory"
	z80_lib "nebula-go/pkg/gbc/z80/lib"
)

func pushReturnAddress(mmu memory.MMU, regs *z80_lib.Registers, value uint16) error {
	sp := regs.SP.Get() - 2

	regs.SP.Set(sp)
	return mmu.WriteDByte(sp, value)
}

func popReturnAddress(mmu memory.MMU, regs *z80_lib.Registers) (uint16, error) {
	sp := regs.SP.Get()

	value, err := mmu.ReadDByte(sp)
	if err != nil {
		return 0, err
	}
	regs.SP.Set(sp + 2)
	return value, nil
}
