package misc

import (
	"nebula-go/pkg/common/bitwise"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
)

func (f *Factory) RLA() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		a := f.regs.A.Get()
		cy := f.regs.F.CY.Get()

		// NOTE: inconsistency here between the two CPU opcodes documents, using official z80 one: ZF is not affected
		//  in official documentation, but is reset in the gameboy opcodes table document.
		f.regs.F.Set(registers.FlagsCleared)
		f.regs.F.CY.Set(bitwise.HighBit8(a))

		f.regs.A.Set((f.regs.A.Get() << 1) | cy)

		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
