package load

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) ConstToByte(reg registerslib.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		value, err := f.mmu.ReadByte(f.regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}
		reg.Set(value)
		return opcodeslib.OpcodeSuccess(2, 8)
	}
}

func (f *Factory) ConstToDByte(reg registerslib.DByte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		value, err := f.mmu.ReadDByte(f.regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}
		reg.Set(value)
		return opcodeslib.OpcodeSuccess(3, 12)
	}
}
