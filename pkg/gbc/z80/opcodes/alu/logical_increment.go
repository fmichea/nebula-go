package alu

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func incrementByteImpl(value uint8, regs *registers.Registers) uint8 {
	value++

	regs.F.ZF.SetBool(value == 0)
	regs.F.HC.SetBool(value&0xF == 0)
	regs.F.NE.SetBool(false)

	return value
}

func (f *Factory) IncrementByte(reg registerslib.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		reg.Set(incrementByteImpl(reg.Get(), f.regs))
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}

func (f *Factory) IncrementDByte(reg registerslib.DByte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		reg.Set(reg.Get() + 1)
		return opcodeslib.OpcodeSuccess(1, 8)
	}
}

func (f *Factory) IncrementHLPtr() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		value, err := f.mmu.ReadByte(f.regs.HL.Get())
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		value = incrementByteImpl(value, f.regs)

		if err := f.mmu.WriteByte(f.regs.HL.Get(), value); err != nil {
			return opcodeslib.OpcodeError(err)
		}

		return opcodeslib.OpcodeSuccess(1, 12)
	}
}
