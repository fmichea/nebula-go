package registers

import (
	"nebula-go/pkg/common/bitwise"
	"nebula-go/pkg/common/frontends"
)

const (
	_paletteCount     = 8
	_colorsPerPalette = 4

	_defaultColor = 0x7FFF
)

type CGBPaletteReg struct {
	indexReg *CGBPaletteIndexReg

	data [][]uint16
}

func NewCGBPaletteReg(index *CGBPaletteIndexReg) *CGBPaletteReg {
	data := make([][]uint16, _paletteCount)
	for idx, _ := range data {
		data[idx] = make([]uint16, _colorsPerPalette)
		for idx2, _ := range data[idx] {
			data[idx][idx2] = _defaultColor
		}
	}

	return &CGBPaletteReg{
		indexReg: index,
		data:     data,
	}
}

func (r *CGBPaletteReg) Get() uint8 {
	paletteIdx, colorIdx, byteIdx := computeByteIndex(r.indexReg.Index.Get())

	v16 := r.data[paletteIdx][colorIdx]
	if byteIdx == 1 {
		return uint8(v16 >> 8)
	}
	return uint8(v16 & 0xFF)
}

func (r *CGBPaletteReg) Set(value uint8) {
	indexValue := r.indexReg.Index.Get()
	paletteIdx, colorIdx, byteIdx := computeByteIndex(indexValue)

	v16 := r.data[paletteIdx][colorIdx]
	if byteIdx == 1 {
		v16 = (uint16(value&0x7F) << 8) | (v16 & 0xFF)
	} else {
		v16 = (v16 & 0xFF00) | uint16(value)
	}
	r.data[paletteIdx][colorIdx] = v16

	if r.indexReg.AutoIncrement.GetBool() {
		r.indexReg.Index.Set(indexValue + 1)
	}
}

func (r *CGBPaletteReg) SetNoMask(value uint8) {
	r.Set(value)
}

func (r *CGBPaletteReg) GetColor(paletteIdx, colorIdx uint8) frontends.Pixel {
	v16 := r.data[paletteIdx][colorIdx]

	return frontends.NewPixel(
		correctColor(bitwise.GetBits16(v16, 0, 0x1F)),
		correctColor(bitwise.GetBits16(v16, 5, 0x1F)),
		correctColor(bitwise.GetBits16(v16, 10, 0x1F)),
	)
}

func computeByteIndex(indexValue uint8) (uint8, uint8, uint8) {
	byteIdx := indexValue % 2
	indexValue /= 2

	colorIdx := indexValue % _colorsPerPalette
	paletteIdx := indexValue / _colorsPerPalette

	return paletteIdx, colorIdx, byteIdx
}

func correctColor(color uint16) uint8 {
	return uint8((color * 0xFF) / 0x1F)
}
