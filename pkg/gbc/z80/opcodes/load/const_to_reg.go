package load

import (
	"nebula-go/pkg/gbc/memory/registers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) ConstToByte(reg registers.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		value, err := f.mmu.ReadByte(f.regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}
		reg.Set(value)
		return opcodeslib.OpcodeSuccess(2, 8)
	}
}

func (f *Factory) ConstToDByte(reg registers.DByte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		value, err := f.mmu.ReadDByte(f.regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}
		reg.Set(value)
		return opcodeslib.OpcodeSuccess(3, 12)
	}
}
