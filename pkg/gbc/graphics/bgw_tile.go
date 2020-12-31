package graphics

import (
	"errors"

	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/memory/registers"
)

const (
	_tilesByLine  = 32
	_bytesPerLine = 2 // bytes
)

func translateXYToAddress(tmds *registers.TMDSFlag, x, y int16) uint16 {
	// Base address depends on the tile map display select.
	base := tmds.BaseAddress()
	// Each byte represents a 8x8 area, so we need to translate the x and y values to that.
	base += uint16(y/8)*_tilesByLine + uint16(x/8)
	// We now have the base address for the background tile for given coordinates.
	return base
}

func getTileBaseAddressAndNumber(mmu memory.MMU, tmds *registers.TMDSFlag, x, y int16) (uint16, uint16, uint8, error) {
	base := translateXYToAddress(tmds, x, y)
	// BG Map Tile Numbers is available directly from base.
	tileNumber, err := mmu.ReadByte(base)
	if err != nil {
		return 0, 0, 0, err
	}
	// Tile address is then resolved using the BG & Window Tile Data Selector.
	tds := mmu.Registers().LCDC.BGWTDS
	return base, tds.BaseAddress(), tds.AdjustTileAddress(tileNumber), err
}

func newBGWTileNonCGB(mmu memory.MMU, tmds *registers.TMDSFlag, x, y int16) (Tile, error) {
	_, tileAddress, tileNumber, err := getTileBaseAddressAndNumber(mmu, tmds, x, y)
	if err != nil {
		return nil, err
	}

	tile := NewDMGBGWTile(mmu, tileTopLeftCorner(x), tileTopLeftCorner(y))

	if err := tile.LoadLineData(mmu, tileAddress, tileNumber, y); err != nil {
		return nil, err
	}

	return tile, nil
}

func wrapInVRAMBankChange(mmu memory.MMU, bank uint8, fn func() error) error {
	vbk := mmu.Registers().VBK

	originalVBK, err := vbk.Get()
	if err != nil {
		return err
	}

	if originalVBK == bank {
		return fn()
	}

	if err := vbk.Set(bank); err != nil {
		return err
	}

	if err := fn(); err != nil {
		return err
	}

	if err := vbk.Set(originalVBK); err != nil {
		return err
	}

	return nil
}

func loadBGWAttributes(mmu memory.MMU, addr uint16) (*TileAttributes, error) {
	var attrs *TileAttributes

	err := wrapInVRAMBankChange(mmu, 1, func() error {
		value, err := mmu.ReadByte(addr)
		if err != nil {
			return err
		}

		attrs = NewTileAttributes(value)
		return nil
	})

	return attrs, err
}

func loadCGBTileDataForLine(mmu memory.MMU, tileAddress uint16, tileNumber uint8, attrs *TileAttributes, tile Tile, y int16) error {
	return wrapInVRAMBankChange(mmu, attrs.VRAMBank.Get(), func() (err error) {
		return tile.LoadLineData(mmu, tileAddress, tileNumber, y)
	})
}

func newBGWTileCGB(mmu memory.MMU, tmds *registers.TMDSFlag, x, y int16) (Tile, error) {
	base, tileAddress, tileNumber, err := getTileBaseAddressAndNumber(mmu, tmds, x, y)
	if err != nil {
		return nil, err
	}

	attrs, err := loadBGWAttributes(mmu, base)
	if err != nil {
		return nil, err
	}

	tile := NewCGBBGWTile(mmu, tileTopLeftCorner(x), tileTopLeftCorner(y), attrs)

	if err := loadCGBTileDataForLine(mmu, tileAddress, tileNumber, attrs, tile, y); err != nil {
		return nil, err
	}

	return tile, nil
}

var (
	ErrBGWXCoordinateInvalid = errors.New("received invalid X coordinate")
	ErrBGWYCoordinateInvalid = errors.New("received invalid Y coordinate")
)

func newBGWTile(mmu memory.MMU, tmds *registers.TMDSFlag, x, y int16) (Tile, error) {
	if x < 0 || 255 < x {
		return nil, ErrBGWXCoordinateInvalid
	}

	if y < 0 || 255 < y {
		return nil, ErrBGWYCoordinateInvalid
	}

	if mmu.Cartridge().IsCGB() {
		return newBGWTileCGB(mmu, tmds, x, y)
	}
	return newBGWTileNonCGB(mmu, tmds, x, y)
}

func NewWindowTile(mmu memory.MMU, x, y int16) (Tile, error) {
	return newBGWTile(mmu, mmu.Registers().LCDC.WTMDS, x, y)
}

func NewBackgroundTile(mmu memory.MMU, x, y int16) (Tile, error) {
	return newBGWTile(mmu, mmu.Registers().LCDC.BGTMDS, x, y)
}

func tileTopLeftCorner(val int16) int16 {
	return val - (val % 8)
}
