package graphics

import (
	"nebula-go/mocks/pkg/gbc/memorymocks"

	"github.com/stretchr/testify/assert"

	"nebula-go/pkg/common/frontends"
	graphicslib "nebula-go/pkg/gbc/graphics/lib"

	"testing"

	"github.com/stretchr/testify/require"
)

// FIXME: this is way too much logic for a test, need to be refactored.
func newTestTile(t *testing.T, mockMMU *memorymocks.MockMMU, x, y int16, data []uint8, hFlip, vFlip, priority bool) Tile {
	paletteFN := func(colorID uint8) frontends.Pixel {
		return frontends.NonePixel
	}

	value := uint8(0)
	if hFlip {
		value |= 1 << 5
	}
	if vFlip {
		value |= 1 << 6
	}
	if priority {
		value |= 1 << 7
	}

	tile := NewTile(x, y, _bgwTileHeight, _objTileContainerWidth, paletteFN, NewTileAttributes(value))

	data16 := make([]uint16, len(data)/2)
	for idx, value := range data {
		v16 := uint16(value)
		if idx%2 == 1 {
			v16 <<= 8
		}
		data16[idx/2] |= v16
	}

	for idx := range data16 {
		dataIdx := idx
		if vFlip {
			dataIdx = len(data16) - 1 - idx
		}

		mockMMU.EXPECT().ReadDByte(uint16(0xF2D0+(2*dataIdx))).Return(data16[dataIdx], nil)
		assert.NoError(t, tile.LoadLineData(mockMMU, 0xF000, 45, y+int16(idx)))
	}

	return tile
}

func buildColumnsFromRows(values [][]uint8) [][]uint8 {
	result := make([][]uint8, 8)

	for _, rowValues := range values {
		for idx, value := range rowValues {
			result[idx] = append(result[idx], value)
		}
	}
	return result
}

func buildReversed(values [][]uint8) [][]uint8 {
	result := make([][]uint8, 8)

	for rowIdx, rowValues := range values {
		result[rowIdx] = make([]uint8, 8)
		for idx, value := range rowValues {
			result[rowIdx][7-idx] = value
		}
	}
	return result
}

func testColumn(t *testing.T, tile Tile, x, yStart int16, expected []uint8) {
	var values []uint8

	for yOffset := int16(0); yOffset < 8; yOffset++ {
		xStart, _, pixels := tile.Colors(yStart + yOffset)
		values = append(values, pixels[x-xStart].ColorID)
	}

	require.Equal(t, expected, values)
}

func testRow(t *testing.T, tile Tile, xStart, y int16, expected []uint8) {
	var values []uint8

	realXStart, _, pixels := tile.Colors(y)
	require.Equal(t, xStart, realXStart)

	for _, pixel := range pixels {
		values = append(values, pixel.ColorID)
	}
	require.Equal(t, expected, values)
}

func (s *unitTestSuite) TestTile_colorID() {
	// The tile should look like this (0 = " ", 1 = ".", 2 = "*", 3 = "#"):
	//    +--------+
	//    |   .*.  |
	//    |   .**. |
	//    |....*#*.|
	//    |*****##*|
	//    |.....  .|
	//    |****. .*|
	//    |###*..*#|
	//    |###*.*##|
	//    +--------+

	tileDataRows := [][]uint8{
		{0, 0, 0, 1, 2, 1, 0, 0},
		{0, 0, 0, 1, 2, 2, 1, 0},
		{1, 1, 1, 1, 2, 3, 2, 1},
		{2, 2, 2, 2, 2, 3, 3, 2},
		{1, 1, 1, 1, 1, 0, 0, 1},
		{2, 2, 2, 2, 1, 0, 1, 2},
		{3, 3, 3, 2, 1, 1, 2, 3},
		{3, 3, 3, 2, 1, 2, 3, 3},
	}

	tileDataRowsReversed := buildReversed(tileDataRows)

	tileDataColumn := buildColumnsFromRows(tileDataRows)
	tileDataColumnReversed := buildReversed(tileDataColumn)

	tileData := []uint8{
		// row 1: 00012100
		0b00010100,
		0b00001000,

		// row 2: 00012210
		0b00010010,
		0b00001100,
		// row 3: 11112321
		0b11110101,
		0b00001110,
		// row 4: 22222332
		0b00000110,
		0b11111111,
		// row 5: 11111001
		0b11111001,
		0b00000000,
		// row 6: 22221012
		0b00001010,
		0b11110001,
		// row 7: 33321123
		0b11101101,
		0b11110011,
		// row 8: 33321233
		0b11101011,
		0b11110111,
	}

	s.Run("simple case, already aligned", func() {
		t := s.T()

		xStart, yStart := int16(0), int16(0)

		tile := newTestTile(t, s.tctx.mockMMU, xStart, yStart, tileData, false, false, false)

		testRow(t, tile, xStart, yStart+0, tileDataRows[0])
		testRow(t, tile, xStart, yStart+2, tileDataRows[2])
		testRow(t, tile, xStart, yStart+4, tileDataRows[4])
		testRow(t, tile, xStart, yStart+6, tileDataRows[6])

		testColumn(t, tile, xStart+0, yStart, tileDataColumn[0])
		testColumn(t, tile, xStart+2, yStart, tileDataColumn[2])
		testColumn(t, tile, xStart+4, yStart, tileDataColumn[4])
		testColumn(t, tile, xStart+6, yStart, tileDataColumn[6])
	})

	s.Run("the tile aligned to 16x32", func() {
		t := s.T()

		xStart, yStart := int16(16), int16(32)

		tile := newTestTile(t, s.tctx.mockMMU, xStart, yStart, tileData, false, false, false)

		testRow(t, tile, xStart, yStart+0, tileDataRows[0])
		testRow(t, tile, xStart, yStart+2, tileDataRows[2])
		testRow(t, tile, xStart, yStart+4, tileDataRows[4])
		testRow(t, tile, xStart, yStart+6, tileDataRows[6])

		testColumn(t, tile, xStart+0, yStart, tileDataColumn[0])
		testColumn(t, tile, xStart+2, yStart, tileDataColumn[2])
		testColumn(t, tile, xStart+4, yStart, tileDataColumn[4])
		testColumn(t, tile, xStart+6, yStart, tileDataColumn[6])
	})

	s.Run("flipped horizontally only", func() {
		t := s.T()

		xStart, yStart := int16(64), int16(64)

		tile := newTestTile(t, s.tctx.mockMMU, xStart, yStart, tileData, true, false, false)

		testRow(t, tile, xStart, yStart+0, tileDataRowsReversed[0])
		testRow(t, tile, xStart, yStart+2, tileDataRowsReversed[2])
		testRow(t, tile, xStart, yStart+4, tileDataRowsReversed[4])
		testRow(t, tile, xStart, yStart+6, tileDataRowsReversed[6])

		testColumn(t, tile, xStart+0, yStart, tileDataColumn[7])
		testColumn(t, tile, xStart+2, yStart, tileDataColumn[5])
		testColumn(t, tile, xStart+4, yStart, tileDataColumn[3])
		testColumn(t, tile, xStart+6, yStart, tileDataColumn[1])
	})

	s.Run("flipped vertically only", func() {
		t := s.T()

		xStart, yStart := int16(23), int16(95)

		tile := newTestTile(t, s.tctx.mockMMU, xStart, yStart, tileData, false, true, false)

		testRow(t, tile, xStart, yStart+0, tileDataRows[7])
		testRow(t, tile, xStart, yStart+2, tileDataRows[5])
		testRow(t, tile, xStart, yStart+4, tileDataRows[3])
		testRow(t, tile, xStart, yStart+6, tileDataRows[1])

		testColumn(t, tile, xStart+0, yStart, tileDataColumnReversed[0])
		testColumn(t, tile, xStart+2, yStart, tileDataColumnReversed[2])
		testColumn(t, tile, xStart+4, yStart, tileDataColumnReversed[4])
		testColumn(t, tile, xStart+6, yStart, tileDataColumnReversed[6])
	})

	s.Run("flipped both horizontally and vertically", func() {
		t := s.T()

		xStart, yStart := int16(13), int16(59)

		tile := newTestTile(t, s.tctx.mockMMU, xStart, yStart, tileData, true, true, false)

		testRow(t, tile, xStart, yStart+0, tileDataRowsReversed[7])
		testRow(t, tile, xStart, yStart+2, tileDataRowsReversed[5])
		testRow(t, tile, xStart, yStart+4, tileDataRowsReversed[3])
		testRow(t, tile, xStart, yStart+6, tileDataRowsReversed[1])

		testColumn(t, tile, xStart+0, yStart, tileDataColumnReversed[7])
		testColumn(t, tile, xStart+2, yStart, tileDataColumnReversed[5])
		testColumn(t, tile, xStart+4, yStart, tileDataColumnReversed[3])
		testColumn(t, tile, xStart+6, yStart, tileDataColumnReversed[1])
	})

	s.Run("aligned below 0 is cut out", func() {
		t := s.T()

		xStart, yStart := int16(-5), int16(-2)

		tile := newTestTile(t, s.tctx.mockMMU, xStart, yStart, tileData, false, false, false)

		testRow(t, tile, 0, yStart+2, tileDataRows[2][5:])
		testRow(t, tile, 0, yStart+4, tileDataRows[4][5:])
	})

	s.Run("aligned close to boundary is cut out", func() {
		t := s.T()

		xStart, yStart := graphicslib.Width-3, int16(4)

		tile := newTestTile(t, s.tctx.mockMMU, xStart, yStart, tileData, false, false, false)

		testRow(t, tile, xStart, yStart+2, tileDataRows[2][:3])
		testRow(t, tile, xStart, yStart+4, tileDataRows[4][:3])
	})

	s.Run("palette returns transparent pixel is nil in Colors", func() {
		t := s.T()

		paletteFN := func(colorID uint8) frontends.Pixel {
			return frontends.TransparentPixel
		}

		tile := NewTile(0, 0, _bgwTileHeight, _objTileContainerWidth, paletteFN, nil)

		_, _, pixels := tile.Colors(0)
		assert.Len(t, pixels, 8)
		assert.Nil(t, pixels[0])
	})
}
