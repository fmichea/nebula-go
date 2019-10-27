package load

import opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"

func (f *Factory) AToHighRAM() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		a8, err := f.mmu.ReadByte(f.regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		if err := f.mmu.WriteByte(0xFF00+uint16(a8), f.regs.A.Get()); err != nil {
			return opcodeslib.OpcodeError(err)
		}

		return opcodeslib.OpcodeSuccess(2, 12)
	}
}

func (f *Factory) HighRAMToA() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		a8, err := f.mmu.ReadByte(f.regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		value, err := f.mmu.ReadByte(0xFF00 + uint16(a8))
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		f.regs.A.Set(value)

		return opcodeslib.OpcodeSuccess(2, 12)
	}
}

func (f *Factory) AToCPtrInHighRAM() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		if err := f.mmu.WriteByte(0xFF00+uint16(f.regs.C.Get()), f.regs.A.Get()); err != nil {
			return opcodeslib.OpcodeError(err)
		}

		return opcodeslib.OpcodeSuccess(1, 8)
	}
}

func (f *Factory) CPtrInHighRAMToA() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		value, err := f.mmu.ReadByte(0xFF00 + uint16(f.regs.C.Get()))
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		f.regs.A.Set(value)
		return opcodeslib.OpcodeSuccess(1, 8)
	}
}
