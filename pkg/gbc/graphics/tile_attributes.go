package graphics

import (
	"nebula-go/pkg/common/bitfields"
)

type TileAttributes struct {
	bitfields.ROByte

	BackgroundPriority bitfields.ROFlag
	VRAMBank           bitfields.ROBitProxy // CGB Only
	DMGPalette         bitfields.ROBitProxy // Used only for OBJ.
	VFlip              bitfields.ROFlag     // CGB Only
	HFlip              bitfields.ROFlag     // CGB Only
	CGBPalette         bitfields.ROBitProxy
}

var tileAttributesCache = make([]*TileAttributes, 0x100)

func NewTileAttributes(value uint8) *TileAttributes {
	obj := tileAttributesCache[value]
	if obj != nil {
		return obj
	}

	reg := bitfields.NewByte(value)

	result := &TileAttributes{
		ROByte: reg,

		BackgroundPriority: bitfields.NewFlag(reg, 7),
		VFlip:              bitfields.NewFlag(reg, 6),
		HFlip:              bitfields.NewFlag(reg, 5),
		DMGPalette:         bitfields.NewBitProxy(reg, 4, 0x1),
		VRAMBank:           bitfields.NewBitProxy(reg, 3, 0x1),
		CGBPalette:         bitfields.NewBitProxy(reg, 0, 0x7),
	}
	tileAttributesCache[value] = result
	return result
}
