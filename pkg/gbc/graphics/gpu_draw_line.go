package graphics

import (
	"errors"
	"nebula-go/pkg/common/frontends"
	graphicslib "nebula-go/pkg/gbc/graphics/lib"
)

func (g *gpu) buildPixelsLayer() []*Pixel {
	return make([]*Pixel, graphicslib.Width)
}

func (g *gpu) getPixelsForLine(ly int16) ([]frontends.Pixel, error) {
	backgroundPixels, err := g.getBackgroundPixels(ly)
	if err != nil {
		return nil, err
	}

	windowPixels, err := g.getWindowPixels(ly)
	if err != nil {
		return nil, err
	}

	objectPixels, err := g.getObjectPixels(ly)
	if err != nil {
		return nil, err
	}

	result := make([]frontends.Pixel, graphicslib.Width)
	for idx := int16(0); idx < graphicslib.Width; idx++ {
		result[idx] = g.selectPixelColor(
			backgroundPixels[idx],
			windowPixels[idx],
			objectPixels[idx],
		)
	}
	return result, nil
}

func (g *gpu) wrapBackgroundAround(value int16) int16 {
	if _bgwTileContainerSize <= value {
		return value - _bgwTileContainerSize
	}
	return value
}

func (g *gpu) shouldRenderBackground() bool {
	return g.cr.IsCGB() || g.mmuRegs.LCDC.BGD.GetBool()
}

func (g *gpu) getBackgroundPixels(ly int16) ([]*Pixel, error) {
	resultPixels := g.buildPixelsLayer()

	if !g.shouldRenderBackground() {
		return resultPixels, nil
	}

	x := int16(g.mmuRegs.SCX.Get())
	backgroundY := g.wrapBackgroundAround(ly + int16(g.mmuRegs.SCY.Get()))

	for pixelsIdx := int16(0); pixelsIdx < graphicslib.Width; {
		wrappedX := g.wrapBackgroundAround(x)

		bgTile, err := NewBackgroundTile(g.mmu, wrappedX, backgroundY)
		if err != nil {
			return nil, err
		}

		startX, endX, tilePixels := bgTile.Colors(backgroundY)
		if len(tilePixels) == 0 {
			return nil, errors.New("background tile returned no pixels")
		}

		if startX < x {
			tilePixels = tilePixels[wrappedX-startX:]
		}

		for _, p := range tilePixels {
			if graphicslib.Width <= pixelsIdx {
				break
			}
			resultPixels[pixelsIdx] = p
			pixelsIdx++
		}

		x = endX
	}

	return resultPixels, nil
}

func (g *gpu) getWindowCoordinates() (int16, int16) {
	windowX := int16(g.mmuRegs.WX.Get())
	windowY := int16(g.mmuRegs.WY.Get())

	if windowX < 7 {
		return 0, windowY
	}
	return windowX - 7, windowY
}

func (g *gpu) getWindowPixels(ly int16) ([]*Pixel, error) {
	resultPixels := g.buildPixelsLayer()

	if !g.mmuRegs.LCDC.WDE.GetBool() {
		return resultPixels, nil
	}

	windowX, windowY := g.getWindowCoordinates()
	if ly < windowY {
		return resultPixels, nil
	}

	offsetY := ly - windowY

	for x := windowX; x < graphicslib.Width; {
		windowTile, err := NewWindowTile(g.mmu, x, offsetY)
		if err != nil {
			return nil, err
		}

		startX, endX, tilePixels := windowTile.Colors(offsetY)
		if len(tilePixels) == 0 {
			return nil, errors.New("window tile returnd no pixels")
		}

		if startX < x {
			tilePixels = tilePixels[x-startX:]
		}

		for idx, p := range tilePixels {
			pixelsIdx := x + int16(idx)
			if graphicslib.Width <= pixelsIdx {
				break
			}
			resultPixels[pixelsIdx] = p
		}

		x = endX
	}

	return resultPixels, nil
}

func (g *gpu) getObjectPixels(ly int16) ([]*Pixel, error) {
	resultPixels := g.buildPixelsLayer()

	if !g.mmuRegs.LCDC.OBJSDE.GetBool() {
		return resultPixels, nil
	}

	sprites, err := g.loadSpritesForLine(ly)
	if err != nil {
		return nil, err
	}

	// Left most sprites are over right most sprites. the loader returns them in
	// order left to right, so we need to iterate the slice returned in reverse
	// order.
	for idx := len(sprites) - 1; 0 <= idx; idx-- {
		startX, _, tilePixels := sprites[idx].Colors(ly)

		for idx, p := range tilePixels {
			pixelsIdx := startX + int16(idx)
			if graphicslib.Width <= pixelsIdx {
				break
			}

			// ColorID 0 for sprites is always fully transparent.
			if p != nil && p.ColorID != 0 {
				resultPixels[pixelsIdx] = p
			}
		}
	}

	return resultPixels, nil
}

func (g *gpu) bgwTileHasPriority(bgPixel, obj *Pixel) bool {
	return (g.mmuRegs.LCDC.BGD.GetBool() && bgPixel.BackgroundPriority) ||
		(bgPixel.ColorID != 0 && obj.BackgroundPriority)
}

func (g *gpu) selectPixelColor(background, window, obj *Pixel) frontends.Pixel {
	bgPixel := background
	if window != nil {
		bgPixel = window
	}

	if obj != nil && (bgPixel == nil || !g.bgwTileHasPriority(bgPixel, obj)) {
		return obj.Color
	} else if bgPixel != nil {
		return bgPixel.Color
	}
	return frontends.White
}

func (g *gpu) drawCurrentLine() error {
	ly := g.mmuRegs.LY.Get()

	pixels, err := g.getPixelsForLine(int16(ly))
	if err != nil {
		return err
	}
	return g.display.DrawLine(uint(ly), pixels)
}
