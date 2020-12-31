package load

import (
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) PushDByte(reg registerslib.DByte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		sp := f.regs.SP.Get() - 2

		f.regs.SP.Set(sp)
		if err := f.mmu.WriteDByte(sp, reg.Get()); err != nil {
			return opcodeslib.OpcodeError(err)
		}
		return opcodeslib.OpcodeSuccess(1, 16)
	}
}

func (f *Factory) PopDByte(reg registerslib.DByte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		sp := f.regs.SP.Get()

		if value, err := f.mmu.ReadDByte(sp); err != nil {
			return opcodeslib.OpcodeError(err)
		} else {
			reg.Set(value)
			f.regs.SP.Set(sp + 2)
			return opcodeslib.OpcodeSuccess(1, 12)
		}
	}
}

func (f *Factory) SPToAddress() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		addr, err := f.mmu.ReadDByte(f.regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		if err := f.mmu.WriteDByte(addr, f.regs.SP.Get()); err != nil {
			return opcodeslib.OpcodeError(err)
		}

		return opcodeslib.OpcodeSuccess(3, 20)
	}
}

// Based on https://stackoverflow.com/questions/5159603/gbz80-how-does-ld-hl-spe-affect-h-and-c-flags the handling for
// carry flags is a bit odd here, since we do sum to a 16 bit value but with only a 8 bit value (my guess), we use
// overflows on bit 4 and 8. So have to implement custom logic here.
func (f *Factory) SPR8ToHL() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		return opcodeslib.SPR8ToDByte(f.mmu, f.regs, f.regs.HL, 12)
	}
}

func (f *Factory) HLToSP() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.regs.SP.Set(f.regs.HL.Get())
		return opcodeslib.OpcodeSuccess(1, 8)
	}
}
