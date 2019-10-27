package alu

import (
	"nebula-go/pkg/common/bitwise"
	"nebula-go/pkg/gbc/memory/registers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) addConstToA(cst uint8) {
	f.logicalOpSetAAndUpdateFlags(cst, false, func(a uint16, value uint16) uint16 {
		return a + value
	})
}

// add %a, $reg
func (f *Factory) AddByteToA(reg registers.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.addConstToA(reg.Get())
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}

func (f *Factory) AddHLPtrToA() opcodeslib.Opcode {
	return f.buildHLPtrToAFunc(f.addConstToA)
}

func (f *Factory) AddD8ToA() opcodeslib.Opcode {
	return f.buildD8ToAFunc(f.addConstToA)
}

func (f *Factory) AddDByteToHL(reg registers.DByte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		value32 := uint32(reg.Get())
		hl32 := uint32(f.regs.HL.Get())

		halfValue := bitwise.Mask32(hl32, 0xFFF) + bitwise.Mask32(value32, 0xFFF)
		fullValue := hl32 + value32

		result := uint16(bitwise.Mask32(fullValue, 0xFFFF))

		f.regs.HL.Set(result)

		f.regs.F.NE.SetBool(false)
		f.regs.F.HC.SetBool(bitwise.InverseMask32(halfValue, 0xFFF) != 0)
		f.regs.F.CY.SetBool(bitwise.InverseMask32(fullValue, 0xFFFF) != 0)

		return opcodeslib.OpcodeSuccess(1, 8)
	}
}

func (f *Factory) AddR8ToSP() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		d8, err := f.mmu.ReadByte(f.regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		f.regs.SP.Set(opcodeslib.AddRelativeConst(f.regs.SP.Get(), d8))
		return opcodeslib.OpcodeSuccess(2, 16)
	}
}
