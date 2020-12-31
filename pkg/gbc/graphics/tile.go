package graphics

import (
	"nebula-go/pkg/common/bitwise"
	"nebula-go/pkg/common/frontends"
	graphicslib "nebula-go/pkg/gbc/graphics/lib"
	"nebula-go/pkg/gbc/memory"
)

const (
	_tileWidth             int16 = 8                 // pixels
	_tileHeight            int16 = 8                 // pixels
	_tileHeightDoubled           = 2 * _tileHeight   // pixels
	_bgwTileHeight         int16 = 8                 // pixels
	_bgwTileContainerSize  int16 = 256               // pixels
	_objTileContainerWidth       = graphicslib.Width // pixels
)

type Pixel struct {
	BackgroundPriority bool
	ColorID            uint8
	Color              frontends.Pixel
}

type PaletteFN func(colorID uint8) frontends.Pixel

type Tile interface {
	LoadLineData(mmu memory.MMU, tileAddress uint16, tileNumber uint8, y int16) error
	Colors(y int16) (startX int16, endX int16, pixels []*Pixel)
}

type tile struct {
	data []uint16

	hFlip              bool
	vFlip              bool
	backgroundPriority bool

	baseX int16
	baseY int16

	height         int16
	containerWidth int16

	paletteFN PaletteFN
}

func NewTile(baseX, baseY, height, containerWidth int16, paletteFN PaletteFN, attrs *TileAttributes) *tile {
	t := &tile{
		data: make([]uint16, height),

		hFlip:              false,
		vFlip:              false,
		backgroundPriority: false,

		baseX: baseX,
		baseY: baseY,

		height:         height,
		containerWidth: containerWidth,

		paletteFN: paletteFN,
	}

	if attrs != nil {
		t.hFlip = attrs.HFlip.GetBool()
		t.vFlip = attrs.VFlip.GetBool()
		t.backgroundPriority = attrs.BackgroundPriority.GetBool()
	}

	return t
}

func (t *tile) LoadLineData(mmu memory.MMU, baseAddress uint16, tileNumber uint8, y int16) error {
	if t.height == _tileHeightDoubled {
		tileNumber = bitwise.InverseMask8(tileNumber, 0x1)
	}

	realY := t.getLineOffset(y)

	// Base address of the tiles depends on input. The tile data here is 8x8
	// in every case. As far as the tileNumber goes, all tiles are 8 pixels
	// high, 2 bytes per line, so we offset by that. Again, there is 2 bytes
	// per line, so we offset that again for the offset of the y line within
	// the tile.
	addr := baseAddress
	addr += 2 * uint16(_tileHeight) * uint16(tileNumber)
	addr += 2 * uint16(realY)

	d16, err := mmu.ReadDByte(addr)
	if err != nil {
		return err
	}

	t.data[realY] = d16
	return nil
}

func (t *tile) Colors(y int16) (int16, int16, []*Pixel) {
	pixels := make([]*Pixel, 0, _tileWidth)

	// We start at the first X coordinate and will go through the tile's whole width
	// (8 pixels). If the tile starts outside of screen (negative X), we adjust it to
	// first pixel, and reduce its width.
	startX, width := t.baseX, _tileWidth
	if startX < 0 {
		width += startX
		startX = 0
	}
	endX := startX + width

	// We iterate over the width of the tile, or up to the width of the container the
	// tile is in, and add those pixels to the result. Background tiles are in a 255x255,
	// window and sprites are in the screen.
	for x := startX; x < endX && x < t.containerWidth; x++ {
		pixels = append(pixels, t.colorID(x, y))
	}
	return startX, startX + int16(len(pixels)), pixels
}

func (t *tile) getLineOffset(y int16) int16 {
	return t.adjustValueToBase(t.baseY, y, t.height, t.vFlip)
}

func (t *tile) colorID(x, y int16) *Pixel {
	realX := uint16(t.adjustValueToBase(t.baseX, x, _tileWidth, t.hFlip))
	realY := t.getLineOffset(y)

	colorID16 := bitwise.GetBit16(t.data[realY], 7-realX)
	colorID16 |= bitwise.GetBit16(t.data[realY], 15-realX) << 1

	colorID := uint8(colorID16)

	color := t.paletteFN(colorID)
	if color == frontends.TransparentPixel {
		return nil
	}

	return &Pixel{
		BackgroundPriority: t.backgroundPriority,
		ColorID:            colorID,
		Color:              color,
	}
}

func (t *tile) adjustValueToBase(base, value, size int16, flip bool) int16 {
	realValue := value - base
	if flip {
		realValue = (size - 1) - realValue
	}
	return realValue
}
