package alu

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) DecrementByte(reg registerslib.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		value := reg.Get() - 1

		reg.Set(value)

		f.regs.F.NE.SetBool(true)
		f.regs.F.ZF.SetBool(value == 0)
		f.regs.F.HC.SetBool(value&0x0F == 0x0F)

		return opcodeslib.OpcodeSuccess(1, 4)
	}
}

func (f *Factory) DecrementDByte(reg registerslib.DByte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		reg.Set(reg.Get() - 1)
		return opcodeslib.OpcodeSuccess(1, 8)
	}
}

func (f *Factory) DecrementHLPtr() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		mmu, regs := f.mmu, f.regs

		value, err := mmu.ReadByte(regs.HL.Get())
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		value--

		if err := mmu.WriteByte(regs.HL.Get(), value); err != nil {
			return opcodeslib.OpcodeError(err)
		}

		regs.F.NE.SetBool(true)
		regs.F.HC.SetBool((value & 0x0F) == 0x0F)
		regs.F.ZF.SetBool(value == 0)

		return opcodeslib.OpcodeSuccess(1, 12)
	}
}
