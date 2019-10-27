package cb

import opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"

func (f *Factory) generateFunctionsBlock(byteFn cbbytefunc, hlPtrFn cbhlptrfunc) []cbopcode {
	return []cbopcode{
		byteFn(f.regs.B),
		byteFn(f.regs.C),
		byteFn(f.regs.D),
		byteFn(f.regs.E),
		byteFn(f.regs.H),
		byteFn(f.regs.L),
		hlPtrFn(),
		byteFn(f.regs.A),
	}
}

func (f *Factory) concatOpcodes(lsts [][]cbopcode) []cbopcode {
	var result []cbopcode

	for _, lst := range lsts {
		result = append(result, lst...)
	}
	return result
}

func (f *Factory) CB() opcodeslib.Opcode {
	opcodeBlocks := [][]cbopcode{
		f.generateFunctionsBlock(f.RLCByte, f.RLCHLPtr),
		f.generateFunctionsBlock(f.RRCByte, f.RRCHLPtr),
		f.generateFunctionsBlock(f.RLByte, f.RLHLPtr),
		f.generateFunctionsBlock(f.RRByte, f.RRHLPtr),
		f.generateFunctionsBlock(f.SLAByte, f.SLAHLPtr),
		f.generateFunctionsBlock(f.SRAByte, f.SRAHLPtr),
		f.generateFunctionsBlock(f.SwapByte, f.SwapHLPtr),
		f.generateFunctionsBlock(f.SRLByte, f.SRLHLPtr),
	}

	for _, fnCtx := range []struct {
		byteFnBuilder  func(uint8) cbbytefunc
		hlPtrFnBuilder func(uint8) cbhlptrfunc
	}{
		{byteFnBuilder: f.TestBitInByte, hlPtrFnBuilder: f.TestBitInHLPtr},
		{byteFnBuilder: f.ResetBitInByte, hlPtrFnBuilder: f.ResetBitInHLPtr},
		{byteFnBuilder: f.SetBitInByte, hlPtrFnBuilder: f.SetBitInHLPtr},
	} {
		for bit := uint8(0); bit < 8; bit++ {
			byteFn := fnCtx.byteFnBuilder(bit)
			hlPtrFn := fnCtx.hlPtrFnBuilder(bit)
			opcodeBlocks = append(opcodeBlocks, f.generateFunctionsBlock(byteFn, hlPtrFn))
		}
	}

	opcodes := f.concatOpcodes(opcodeBlocks)

	return func() opcodeslib.OpcodeResult {
		value, err := f.mmu.ReadByte(f.regs.PC + 1)
		if err != nil {
			return opcodeslib.OpcodeError(err)
		}

		return opcodes[value]()
	}
}
