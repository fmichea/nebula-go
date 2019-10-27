package load

import (
	"nebula-go/pkg/gbc/memory/registers"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
)

func (f *Factory) loadAToAddress(addr uint16) error {
	if err := f.mmu.WriteByte(addr, f.regs.A.Get()); err != nil {
		return err
	}
	return nil
}

func (f *Factory) aToDBytePtr(reg registers.DByte) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		if err := f.loadAToAddress(reg.Get()); err != nil {
			return opcodeslib.OpcodeError(err)
		}
		return opcodeslib.OpcodeSuccess(1, 8)
	}
}

func (f *Factory) aToHLPtrMod(fn func(val uint16) uint16) opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		hl := f.regs.HL.Get()

		if err := f.loadAToAddress(hl); err != nil {
			return opcodeslib.OpcodeError(err)
		}
		f.regs.HL.Set(fn(hl))
		return opcodeslib.OpcodeSuccess(1, 8)
	}
}

func (f *Factory) AToAddress() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		a16, err := f.mmu.ReadDByte(f.regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		if err := f.loadAToAddress(a16); err != nil {
			return opcodeslib.OpcodeError(err)
		}

		return opcodeslib.OpcodeSuccess(3, 16)
	}
}

func (f *Factory) AToBCPtr() opcodeslib.Opcode {
	return f.aToDBytePtr(f.regs.BC)
}

func (f *Factory) AToDEPtr() opcodeslib.Opcode {
	return f.aToDBytePtr(f.regs.DE)
}

func (f *Factory) AToHLInc() opcodeslib.Opcode {
	return f.aToHLPtrMod(func(val uint16) uint16 {
		return val + 1
	})
}

func (f *Factory) AToHLDec() opcodeslib.Opcode {
	return f.aToHLPtrMod(func(val uint16) uint16 {
		return val - 1
	})
}
