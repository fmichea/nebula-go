package graphics

import (
	"nebula-go/pkg/common/frontends"
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/memory/registers"
)

func NewCGBBGWTile(mmu memory.MMU, x, y int16, attrs *TileAttributes) Tile {
	return newCGBTile(x, y, _bgwTileHeight, _bgwTileContainerSize, mmu.Registers().BGPD, attrs)
}

func NewCGBObjTile(mmu memory.MMU, x, y, height int16, attrs *TileAttributes) Tile {
	return newCGBTile(x, y, height, _objTileContainerWidth, mmu.Registers().OBPD, attrs)
}

func newCGBTile(x, y, height, containerWidth int16, palette *registers.CGBPaletteReg, attrs *TileAttributes) Tile {
	paletteIdx := attrs.CGBPalette.Get()
	paletteFN := func(colorID uint8) frontends.Pixel {
		return palette.GetColor(paletteIdx, colorID)
	}
	return NewTile(x, y, height, containerWidth, paletteFN, attrs)
}
