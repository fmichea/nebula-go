package load

import (
	"nebula-go/pkg/gbc/memory/registers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) loadAddressToA(addr uint16) error {
	value, err := f.mmu.ReadByte(addr)
	if err != nil {
		return err
	}

	f.regs.A.Set(value)
	return nil
}

func (f *Factory) dbytePtrToA(reg registers.DByte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		if err := f.loadAddressToA(reg.Get()); err != nil {
			return opcodeslib.OpcodeError(err)
		}
		return opcodeslib.OpcodeSuccess(1, 8)
	}
}

func (f *Factory) hlPtrModToA(fn func(val uint16) uint16) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		hl := f.regs.HL.Get()

		if err := f.loadAddressToA(hl); err != nil {
			return opcodeslib.OpcodeError(err)
		}
		f.regs.HL.Set(fn(hl))
		return opcodeslib.OpcodeSuccess(1, 8)
	}
}

func (f *Factory) AddressToA() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		a16, err := f.mmu.ReadDByte(f.regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		if err := f.loadAddressToA(a16); err != nil {
			return opcodeslib.OpcodeError(err)
		}

		return opcodeslib.OpcodeSuccess(3, 16)
	}
}

func (f *Factory) BCPtrToA() opcodeslib.Opcode {
	return f.dbytePtrToA(f.regs.BC)
}

func (f *Factory) DEPtrToA() opcodeslib.Opcode {
	return f.dbytePtrToA(f.regs.DE)
}

func (f *Factory) HLIncToA() opcodeslib.Opcode {
	return f.hlPtrModToA(func(val uint16) uint16 {
		return val + 1
	})
}

func (f *Factory) HLDecToA() opcodeslib.Opcode {
	return f.hlPtrModToA(func(val uint16) uint16 {
		return val - 1
	})
}
