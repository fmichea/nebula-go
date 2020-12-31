package graphics

import (
	graphicslib "nebula-go/pkg/gbc/graphics/lib"
	"nebula-go/pkg/gbc/memory/registers"
)

const (
	// LY will increment from 0 to 153, where the values between 144 and 153 are V-Blank.
	_maxLY uint8 = 153
)

func (g *gpu) doHBlank() error {
	ly := g.mmuRegs.LY.Get()

	height8 := uint8(graphicslib.Height)

	coincidence := ly == g.mmuRegs.LYC.Get()
	if coincidence {
		g.mmuRegs.IF.STAT.SetBool(true)
	}
	g.mmuRegs.STAT.CoincidenceFlag.SetBool(coincidence)

	if g.mmuRegs.LCDC.LDE.GetBool() && ly < height8 {
		if err := g.drawCurrentLine(); err != nil {
			return err
		}
	}

	// Until last line is drawn, we keep doing the data transfer cycles. Once
	// last line is drawn, we switch to a V-Blank.
	g.mmuRegs.LY.Set(ly + 1)
	if ly < height8-1 {
		g.mmuRegs.STAT.Mode.SetMode(registers.STATModeDataTransfer1)
	} else {
		g.mmuRegs.STAT.Mode.SetMode(registers.STATModeVBlank)
	}

	if g.cr.IsCGB() {
		return g.mmu.Registers().HDMA5.MaybeDoHDMA()
	}
	return nil
}

func (g *gpu) doVBlank() error {
	ly := g.mmuRegs.LY.Get()

	if ly < _maxLY {
		g.mmuRegs.STAT.Mode.SetMode(registers.STATModeVBlank)
		g.mmuRegs.LY.Set(ly + 1)

		return nil
	} else {
		g.pacer.Wait()

		g.mmuRegs.STAT.Mode.SetMode(registers.STATModeDataTransfer1)
		g.mmuRegs.IF.VBlank.SetBool(true)
		g.mmuRegs.LY.Set(0)

		return g.display.Commit()
	}
}

func (g *gpu) doDataTransfer1() error {
	g.mmuRegs.STAT.Mode.SetMode(registers.STATModeDataTransfer2)
	return nil
}

func (g *gpu) doDataTransfer2() error {
	g.mmuRegs.STAT.Mode.SetMode(registers.STATModeHBlank)
	return nil
}

func (g *gpu) DoCycles(cycles uint16) error {
	timer := g.mmuRegs.STAT.Timer

	timer.Forward(cycles)
	if !timer.Expired() {
		return nil
	}

	switch g.mmuRegs.STAT.Mode.GetMode() {
	case registers.STATModeHBlank:
		return g.doHBlank()

	case registers.STATModeVBlank:
		return g.doVBlank()

	case registers.STATModeDataTransfer1:
		return g.doDataTransfer1()

	case registers.STATModeDataTransfer2:
		return g.doDataTransfer2()
	}

	return nil
}
