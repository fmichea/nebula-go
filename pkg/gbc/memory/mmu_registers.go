package memory

import (
	"nebula-go/pkg/gbc/memory/cartridge"
	"nebula-go/pkg/gbc/memory/lib"
	"nebula-go/pkg/gbc/memory/segments"

	"go.uber.org/multierr"

	"nebula-go/pkg/gbc/memory/registers"
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

type Registers struct {
	Stopped bool

	KEY1 *registers.KEY1Reg // CGB speed control register.

	LCDC *registers.LCDCReg // LCD Control (R/W)
	STAT *registers.STATReg // LCDC Status (R/W)

	SCY registerslib.Byte // Scroll Y (R/W)
	SCX registerslib.Byte // Scroll X (R/W)

	LY  registerslib.Byte // LCDC Y-Coordinate (R) FIXME: writing between 144 and 153 resets this register.
	LYC registerslib.Byte // LY Compare (R/W)

	WY registerslib.Byte // Window Y Position (R/W)
	WX registerslib.Byte // Window X Position minus 7 (R/W)

	BGP  *registers.DMGPaletteReg // BG Palette Data (R/W) - Non CGB Mode Only
	OBP0 *registers.DMGPaletteReg // Object Palette 0 Data (R/W) - Non CGB Mode Only
	OBP1 *registers.DMGPaletteReg // Object Palette 1 Data (R/W) - Non CGB Mode Only

	BGPI *registers.CGBPaletteIndexReg
	BGPD *registers.CGBPaletteReg

	OBPI *registers.CGBPaletteIndexReg
	OBPD *registers.CGBPaletteReg

	VBK  *registers.VBKReg
	SVBK *registers.VBKReg

	IF *registers.InterruptReg
	IE *registers.InterruptReg

	JOYP *registers.JOYPReg

	DIV  *registers.DIVReg
	TIMA *registers.TIMAReg
	TMA  *registers.TMAReg
	TAC  *registers.TACReg

	DMA *registers.DMAReg

	HDMA1 registerslib.Byte
	HDMA2 registerslib.Byte
	HDMA3 registerslib.Byte
	HDMA4 registerslib.Byte
	HDMA5 registers.HDMA5Reg

	maskedRegisters map[string]registerslib.Byte
}

func InitializeRegisters(mmu lib.MemoryIO, cr *cartridge.ROM, vram, wram segments.Segment) (*Registers, error) {
	// Registers mapped into MMU.
	regsMapping := map[string]registerslib.Byte{}
	for _, cfg := range []struct {
		name string
		addr uint16
	}{
		{"SCY", 0xFF42},
		{"SCX", 0xFF43},
		{"LY", 0xFF44},
		{"LYC", 0xFF45},
		{"WY", 0xFF4A},
		{"WX", 0xFF4B},
	} {
		err := mmu.ByteHook(cfg.addr, func(ptr *uint8) (hook lib.Hook, e error) {
			regsMapping[cfg.name] = registerslib.NewByte(ptr, 0)
			hook = registerslib.WrapWithError(regsMapping[cfg.name])
			return
		})
		if err != nil {
			return nil, err
		}
	}

	dmgPaletteMapping := map[string]*registers.DMGPaletteReg{}
	for _, cfg := range []struct {
		name         string
		addr         uint16
		value        uint8
		transparent0 bool
	}{
		{"BGP", 0xFF47, 0xFC, false},
		{"OBP0", 0xFF48, 0xFF, true},
		{"OBP1", 0xFF49, 0xFF, true},
	} {
		err := mmu.ByteHook(cfg.addr, func(ptr *uint8) (hook lib.Hook, e error) {
			dmgPaletteMapping[cfg.name] = registers.NewDMGPaletteReg(ptr, cfg.value, cfg.transparent0)
			hook = registerslib.WrapWithError(dmgPaletteMapping[cfg.name])
			return
		})
		if err != nil {
			return nil, err
		}
	}

	cgbPaletteIndexMapping := map[string]*registers.CGBPaletteIndexReg{}
	for _, cfg := range []struct {
		name string
		addr uint16
	}{
		{"BGPI", 0xFF68},
		{"OBPI", 0xFF6A},
	} {
		err := mmu.ByteHook(cfg.addr, func(ptr *uint8) (hook lib.Hook, e error) {
			cgbPaletteIndexMapping[cfg.name] = registers.NewCGBPaletteIndexReg(ptr)
			hook = registerslib.WrapWithError(cgbPaletteIndexMapping[cfg.name])
			return
		})
		if err != nil {
			return nil, err
		}
	}

	cgbPaletteMapping := map[string]*registers.CGBPaletteReg{}
	for _, cfg := range []struct {
		name      string
		addr      uint16
		indexName string
	}{
		{"BGPD", 0xFF69, "BGPI"},
		{"OBPD", 0xFF6B, "OBPI"},
	} {
		err := mmu.ByteHook(cfg.addr, func(ptr *uint8) (hook lib.Hook, e error) {
			cgbPaletteMapping[cfg.name] = registers.NewCGBPaletteReg(cgbPaletteIndexMapping[cfg.indexName])
			hook = registerslib.WrapWithError(cgbPaletteMapping[cfg.name])
			return
		})
		if err != nil {
			return nil, err
		}
	}

	interruptRegsMapping := map[string]*registers.InterruptReg{}
	for _, cfg := range []struct {
		name string
		addr uint16
	}{
		{"IF", 0xFF0F},
		{"IE", 0xFFFF},
	} {
		err := mmu.ByteHook(cfg.addr, func(ptr *uint8) (hook lib.Hook, e error) {
			interruptRegsMapping[cfg.name] = registers.NewInterruptReg(ptr)
			hook = registerslib.WrapWithError(interruptRegsMapping[cfg.name])
			return
		})
		if err != nil {
			return nil, err
		}
	}

	maskedRegsMapping := map[string]registerslib.Byte{}
	for _, cfg := range []struct {
		name         string
		addr         uint16
		initialValue uint8
		mask         uint8
	}{
		// FIXME: add the sound registers here?

		{"FF6C", 0xFF6C, 0xFE, 0x01},
		{"FF75", 0xFF75, 0x8F, 0x38},
		{"FF76", 0xFF76, 0x00, 0x00},
		{"FF77", 0xFF77, 0x00, 0x00},
	} {
		err := mmu.ByteHook(cfg.addr, func(ptr *uint8) (hook lib.Hook, err error) {
			maskedRegsMapping[cfg.name] = registerslib.NewByteWithMask(ptr, cfg.initialValue, cfg.mask)
			hook = registerslib.WrapWithError(maskedRegsMapping[cfg.name])
			return
		})
		if err != nil {
			return nil, err
		}
	}

	regs := &Registers{
		Stopped: false,

		SCY: regsMapping["SCY"],
		SCX: regsMapping["SCX"],

		LY:  regsMapping["LY"],
		LYC: regsMapping["LYC"],

		WY: regsMapping["WY"],
		WX: regsMapping["WX"],

		BGP:  dmgPaletteMapping["BGP"],
		OBP0: dmgPaletteMapping["OBP0"],
		OBP1: dmgPaletteMapping["OBP1"],

		BGPI: cgbPaletteIndexMapping["BGPI"],
		BGPD: cgbPaletteMapping["BGPD"],

		OBPI: cgbPaletteIndexMapping["OBPI"],
		OBPD: cgbPaletteMapping["OBPD"],

		IF: interruptRegsMapping["IF"],
		IE: interruptRegsMapping["IE"],

		maskedRegisters: maskedRegsMapping,
	}

	err := multierr.Combine(
		mmu.ByteHook(0xFF51, func(ptr *uint8) (hook lib.Hook, err error) {
			regs.HDMA1 = registers.NewHDMA1Reg(ptr)
			hook = registerslib.WrapWithError(regs.HDMA1)
			return
		}),
		mmu.ByteHook(0xFF52, func(ptr *uint8) (hook lib.Hook, err error) {
			regs.HDMA2 = registers.NewHDMA2Reg(ptr)
			hook = registerslib.WrapWithError(regs.HDMA2)
			return
		}),
		mmu.ByteHook(0xFF53, func(ptr *uint8) (hook lib.Hook, err error) {
			regs.HDMA3 = registers.NewHDMA3Reg(ptr)
			hook = registerslib.WrapWithError(regs.HDMA3)
			return
		}),
		mmu.ByteHook(0xFF54, func(ptr *uint8) (hook lib.Hook, err error) {
			regs.HDMA4 = registers.NewHDMA4Reg(ptr)
			hook = registerslib.WrapWithError(regs.HDMA4)
			return
		}),
		// NOTE: This must come after HDMA[1-4]
		mmu.ByteHook(0xFF55, func(ptr *uint8) (hook lib.Hook, err error) {
			regs.HDMA5 = registers.NewHDMA5Reg(ptr, mmu, regs.HDMA1, regs.HDMA2, regs.HDMA3, regs.HDMA4)
			hook = regs.HDMA5
			return
		}),
		mmu.ByteHook(0xFF40, func(ptr *uint8) (hook lib.Hook, e error) {
			regs.LCDC = registers.NewLCDCReg(ptr)
			hook = registerslib.WrapWithError(regs.LCDC)
			return
		}),
		mmu.ByteHook(0xFF41, func(ptr *uint8) (hook lib.Hook, e error) {
			regs.STAT = registers.NewSTATReg(ptr)
			hook = registerslib.WrapWithError(regs.STAT)
			return
		}),
		mmu.ByteHook(0xFF4F, func(ptr *uint8) (hook lib.Hook, e error) {
			// FIXME: add check for HDMA5 activity here.
			regs.VBK, e = registers.NewVBKReg(ptr, cr, vram)
			hook = regs.VBK
			return
		}),
		mmu.ByteHook(0xFF70, func(ptr *uint8) (hook lib.Hook, e error) {
			regs.SVBK, e = registers.NewVBKReg(ptr, cr, wram)
			hook = regs.SVBK
			return
		}),
		mmu.ByteHook(0xFF4D, func(ptr *uint8) (hook lib.Hook, err error) {
			regs.KEY1 = registers.NewKEY1Reg(ptr)
			hook = registerslib.WrapWithError(regs.KEY1)
			return
		}),
		mmu.ByteHook(0xFF00, func(ptr *uint8) (hook lib.Hook, err error) {
			regs.JOYP = registers.NewJOYPReg(ptr)
			hook = registerslib.WrapWithError(regs.JOYP)
			return
		}),
		mmu.ByteHook(0xFF07, func(ptr *uint8) (hook lib.Hook, err error) {
			regs.TAC = registers.NewTACReg(ptr)
			hook = registerslib.WrapWithError(regs.TAC)
			return
		}),
		// NOTE: this must be after TAC.
		mmu.ByteHook(0xFF04, func(ptr *uint8) (hook lib.Hook, err error) {
			regs.DIV = registers.NewDIVReg(ptr, regs.TAC)
			hook = registerslib.WrapWithError(regs.DIV)
			return
		}),
		mmu.ByteHook(0xFF06, func(ptr *uint8) (hook lib.Hook, err error) {
			regs.TMA = registers.NewTMAReg(ptr)
			hook = registerslib.WrapWithError(regs.TMA)
			return
		}),
		// NOTE: this must be after DIV, TMA and IF
		mmu.ByteHook(0xFF05, func(ptr *uint8) (hook lib.Hook, err error) {
			regs.TIMA = registers.NewTIMAReg(ptr, regs.DIV, regs.TMA, regs.IF)
			hook = registerslib.WrapWithError(regs.TIMA)
			return
		}),
		mmu.ByteHook(0xFF46, func(ptr *uint8) (hook lib.Hook, err error) {
			regs.DMA = registers.NewDMAReg(mmu)
			hook = regs.DMA
			return
		}),
	)
	if err != nil {
		return nil, err
	}

	return regs, nil
}
