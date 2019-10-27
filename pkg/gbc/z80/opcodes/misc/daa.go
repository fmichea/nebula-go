package misc

import opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"

// DAA fixes values after a BCD operation on the accumulator. This contains quite a few branches which are explained
// in the code, but you can find a better explanation of what this all does in
// Explanation of this instruction: https://ehaskins.com/2018-01-30%20Z80%20DAA/
func (f *Factory) DAA() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		regs := f.regs

		a16 := uint16(regs.A.Get())

		if regs.F.NE.Get() != 0 {
			if regs.F.HC.GetBool() {
				a16 = (a16 - 6) & 0xFF
			}
			if regs.F.CY.GetBool() {
				a16 -= 0x60
			}
		} else {
			if regs.F.HC.GetBool() || 9 < (a16&0x0F) {
				a16 += 0x06
			}
			if regs.F.CY.GetBool() || 0x9F < a16 {
				a16 += 0x60
			}
		}

		a8 := uint8(a16)

		regs.A.Set(a8)
		regs.F.ZF.SetBool(a8 == 0)
		regs.F.HC.SetBool(false)
		regs.F.CY.SetBool(regs.F.CY.GetBool() || 0xFF < a16)

		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
