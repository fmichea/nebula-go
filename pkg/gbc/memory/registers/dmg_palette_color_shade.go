package registers

import (
	"nebula-go/pkg/common/frontends"
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

var (
	_black     = frontends.Black
	_darkGrey  = frontends.NewPixel(0x55, 0x55, 0x55)
	_lightGrey = frontends.NewPixel(0xAA, 0xAA, 0xAA)
	_white     = frontends.White

	_dmgColors = []frontends.Pixel{_white, _lightGrey, _darkGrey, _black}
)

func NewDMGPaletteColorShade(reg registerslib.Byte, offset uint8) *DMGPaletteColorShade {
	return &DMGPaletteColorShade{
		BitProxy: registerslib.NewBitProxy(reg, offset, 0x3),
	}
}

type DMGPaletteColorShade struct {
	registerslib.BitProxy
}

func (s *DMGPaletteColorShade) GetColor() frontends.Pixel {
	return _dmgColors[s.Get()]
}
