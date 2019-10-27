package alu

import (
	"nebula-go/pkg/common/bitwise"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

type logicalFunc func(a uint16, value uint16) uint16

func (f *Factory) logicalOpUpdateFlags(value uint8, isNeg bool, fn logicalFunc) uint8 {
	a16 := uint16(f.regs.A.Get())
	value16 := uint16(value)

	halfValue := fn(bitwise.Mask16(a16, 0x0F), bitwise.Mask16(value16, 0x0F))
	fullValue := fn(a16, value16)

	result := uint8(bitwise.Mask16(fullValue, 0xFF))

	f.regs.F.ZF.SetBool(result == 0)
	f.regs.F.NE.SetBool(isNeg)
	f.regs.F.HC.SetBool(bitwise.InverseMask16(halfValue, 0xF) != 0)
	f.regs.F.CY.SetBool(bitwise.InverseMask16(fullValue, 0xFF) != 0)

	return result
}

func (f *Factory) logicalOpSetAAndUpdateFlags(value uint8, isNeg bool, fn logicalFunc) {
	f.regs.A.Set(f.logicalOpUpdateFlags(value, isNeg, fn))
}

func (f *Factory) buildHLPtrToAFunc(fn func(cst uint8)) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		mhl, err := f.mmu.ReadByte(f.regs.HL.Get())
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		fn(mhl)
		return opcodeslib.OpcodeSuccess(1, 8)
	}
}

func (f *Factory) buildD8ToAFunc(fn func(cst uint8)) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		d8, err := f.mmu.ReadByte(f.regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		fn(d8)
		return opcodeslib.OpcodeSuccess(2, 8)
	}
}
