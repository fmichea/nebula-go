package cb

import (
	lib2 "nebula-go/pkg/gbc/z80/opcodes/lib"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

func (f *Factory) testBit(value, bit uint8) {
	bitValue := (value >> bit) & 0x1

	f.regs.F.ZF.SetBool(bitValue == 0)
	f.regs.F.NE.SetBool(false)
	f.regs.F.HC.SetBool(true)
}

func (f *Factory) TestBitInByte(bit uint8) cbbytefunc {
	return func(reg registerslib.Byte) cbopcode {
		return func() lib2.OpcodeResult {
			f.testBit(reg.Get(), bit)
			return lib2.OpcodeSuccess(2, 8)
		}
	}
}

func (f *Factory) TestBitInHLPtr(bit uint8) cbhlptrfunc {
	return func() cbopcode {
		return func() lib2.OpcodeResult {
			value, err := f.mmu.ReadByte(f.regs.HL.Get())
			if err != nil {
				return lib2.OpcodeError(err)
			}

			f.testBit(value, bit)

			// NOTE: Instruction timing ROM uses 12 cycles for this instruction, not 16 like
			//  the opcode documentation indicates.
			return lib2.OpcodeSuccess(2, 12)
		}
	}
}
