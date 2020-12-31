package graphics

import (
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/memory/registers"
)

func NewDMGBGWTile(mmu memory.MMU, baseX, baseY int16) Tile {
	palettes := []*registers.DMGPaletteReg{mmu.Registers().BGP}
	return newDMGTile(baseX, baseY, _bgwTileHeight, _bgwTileContainerSize, nil, palettes)
}

func NewDMGObjTile(mmu memory.MMU, baseX, baseY, height int16, attrs *TileAttributes) Tile {
	mmuRegs := mmu.Registers()
	palettes := []*registers.DMGPaletteReg{mmuRegs.OBP0, mmuRegs.OBP1}
	return newDMGTile(baseX, baseY, height, _objTileContainerWidth, attrs, palettes)
}

func newDMGTile(x, y, height, containerWidth int16, attrs *TileAttributes, palettes []*registers.DMGPaletteReg) Tile {
	paletteIdx := uint8(0)
	if attrs != nil {
		paletteIdx = attrs.DMGPalette.Get()
	}
	return NewTile(x, y, height, containerWidth, palettes[paletteIdx].GetColor, attrs)
}
