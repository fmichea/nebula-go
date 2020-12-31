package graphics

import (
	"sort"

	graphicslib "nebula-go/pkg/gbc/graphics/lib"
)

const (
	_spriteTilesBaseAddress uint16 = 0x8000

	_spriteCount                   = 40
	_spriteAttributesByteSize      = 4
	_spriteAttributesTableByteSize = _spriteCount * _spriteAttributesByteSize

	_maxSpritesPerScanLine = 10
)

type sprite struct {
	x int16

	tile Tile
}

func (g *gpu) loadSpritesForLine(y int16) ([]Tile, error) {
	sprites := make([]sprite, 0, _maxSpritesPerScanLine)

	spritesData, err := g.mmu.ReadByteSlice(0xFE00, _spriteAttributesTableByteSize)
	if err != nil {
		return nil, err
	}

	spriteHeight := g.mmu.Registers().LCDC.OBJSS.SpriteHeight()

	for idx, spritesFound := 0, 0; idx < _spriteCount && spritesFound < _maxSpritesPerScanLine; idx++ {
		spriteDataIdx := _spriteAttributesByteSize * idx

		// The sprite Y value is adjusted by 16 pixels (left).
		spriteY := int16(spritesData[spriteDataIdx]) - 16
		spriteEndY := spriteY + spriteHeight - 1

		// We only care about sprites that are visible in the current scan line.
		if !isInBounds(y, spriteY, spriteEndY) {
			continue
		}

		// Even if it is hidden by its X value, this sprite is counted as present for scanline purposes.
		spritesFound++

		// The sprite X value is adjusted by 8 pixels (up).
		spriteX := int16(spritesData[spriteDataIdx+1]) - 8
		spriteEndX := spriteX + _tileWidth

		// Any part of the sprite being present on the line is displayed.
		if !isInBounds(spriteX, 0, graphicslib.Width-1) && !isInBounds(spriteEndX, 1, graphicslib.Width) {
			continue
		}

		tileNumber := spritesData[spriteDataIdx+2]
		attrs := NewTileAttributes(spritesData[spriteDataIdx+3])

		tileFN := NewDMGObjTile
		if g.cr.IsCGB() {
			tileFN = NewCGBObjTile
		}

		tile := tileFN(g.mmu, spriteX, spriteY, spriteHeight, attrs)

		if err := tile.LoadLineData(g.mmu, _spriteTilesBaseAddress, tileNumber, y); err != nil {
			return nil, err
		}

		sprites = append(sprites, sprite{
			x:    spriteX,
			tile: tile,
		})
	}

	// Sprites will be returned in order from left to right (or when first found
	// during scan line on same X) for DMG. For CGB, table ordering is always
	// used.
	if !g.cr.IsCGB() {
		sort.SliceStable(sprites, func(i, j int) bool {
			iX, jX := sprites[i].x, sprites[j].x
			return iX < jX
		})
	}

	result := make([]Tile, 0, _maxSpritesPerScanLine)
	for _, sprite := range sprites {
		result = append(result, sprite.tile)
	}

	return result, nil
}

func isInBounds(value, lowerBound, upperBound int16) bool {
	return lowerBound <= value && value <= upperBound
}
