package load

import (
	"nebula-go/pkg/gbc/memory/registers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

// LoadHLPtrToByte implements the following opcode: ld $reg, (%hl)
func (f *Factory) HLPtrToByte(reg registers.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		value, err := f.mmu.ReadByte(f.regs.HL.Get())
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}
		reg.Set(value)
		return opcodeslib.OpcodeSuccess(1, 8)
	}
}

func (f *Factory) ByteToHLPtr(reg registers.Byte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		if err := f.mmu.WriteByte(f.regs.HL.Get(), reg.Get()); err != nil {
			return opcodeslib.OpcodeError(err)
		}
		return opcodeslib.OpcodeSuccess(1, 8)
	}
}

func (f *Factory) D8ToHLPtr() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		value, err := f.mmu.ReadByte(f.regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		if err := f.mmu.WriteByte(f.regs.HL.Get(), value); err != nil {
			return opcodeslib.OpcodeError(err)
		}

		return opcodeslib.OpcodeSuccess(2, 12)
	}
}
