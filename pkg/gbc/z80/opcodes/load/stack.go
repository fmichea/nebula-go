package load

import (
	"nebula-go/pkg/common/bitwise"
	"nebula-go/pkg/gbc/memory/registers"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) PushDByte(reg registers.DByte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		sp := f.regs.SP.Get() - 2

		f.regs.SP.Set(sp)
		if err := f.mmu.WriteDByte(sp, reg.Get()); err != nil {
			return opcodeslib.OpcodeError(err)
		}
		return opcodeslib.OpcodeSuccess(1, 16)
	}
}

func (f *Factory) PopDByte(reg registers.DByte) opcodeslib.Opcode {
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
		sp := f.regs.SP.Get()

		d8, err := f.mmu.ReadByte(f.regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		fn := opcodeslib.AddRelativeConstForMaskFunc(sp, d8)
		carryFn := func(mask uint16) bool {
			return bitwise.InverseMask16(fn(mask), mask) != 0
		}

		f.regs.HL.Set(fn(0xFFFF))

		f.regs.F.Set(z80lib.FlagsCleared)
		f.regs.F.HC.SetBool(carryFn(0x000F))
		f.regs.F.CY.SetBool(carryFn(0x00FF))

		return opcodeslib.OpcodeSuccess(2, 12)
	}
}

func (f *Factory) HLToSP() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.regs.SP.Set(f.regs.HL.Get())
		return opcodeslib.OpcodeSuccess(1, 8)
	}
}
