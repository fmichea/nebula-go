package cb

import (
	"nebula-go/pkg/gbc/memory/registers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

type cbopcode func() opcodeslib.OpcodeResult
type cbbytefunc func(registers.Byte) cbopcode
type cbhlptrfunc func() cbopcode

func (f *Factory) updateFlagsForCB(value, cy uint8) uint8 {
	f.regs.F.ZF.SetBool(value == 0)
	f.regs.F.NE.SetBool(false)
	f.regs.F.HC.SetBool(false)
	f.regs.F.CY.Set(cy)
	return value
}

func (f *Factory) buildCBOpcodeByte(reg registers.Byte, fn func(uint8) uint8) cbopcode {
	return func() opcodeslib.OpcodeResult {
		reg.Set(fn(reg.Get()))
		return opcodeslib.OpcodeSuccess(2, 8)
	}
}

func (f *Factory) buildCBOpcodeHLPtr(fn func(uint8) uint8) cbopcode {
	return func() opcodeslib.OpcodeResult {
		hl := f.regs.HL.Get()

		value, err := f.mmu.ReadByte(hl)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		newValue := fn(value)

		if err := f.mmu.WriteByte(hl, newValue); err != nil {
			return opcodeslib.OpcodeError(err)
		}

		return opcodeslib.OpcodeSuccess(2, 16)
	}
}
