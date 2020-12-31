package misc

import opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"

func (f *Factory) Stop() opcodeslib.Opcode {
	return func() opcodeslib.OpcodeResult {
		f.mmu.Registers().KEY1.SwitchIfRequested()
		return opcodeslib.OpcodeSuccess(1, 4)
	}
}
