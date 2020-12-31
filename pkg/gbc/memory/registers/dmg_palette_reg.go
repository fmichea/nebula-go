package registers

import (
	"nebula-go/pkg/common/frontends"
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

func NewDMGPaletteReg(ptr *uint8, value uint8, transparent0 bool) *DMGPaletteReg {
	reg := registerslib.NewByte(ptr, value)

	shade0 := NewDMGPaletteColorShade(reg, 0)
	if transparent0 {
		shade0 = nil
	}

	return &DMGPaletteReg{
		Byte: reg,

		Shades: []*DMGPaletteColorShade{
			shade0,
			NewDMGPaletteColorShade(reg, 2),
			NewDMGPaletteColorShade(reg, 4),
			NewDMGPaletteColorShade(reg, 6),
		},
	}
}

type DMGPaletteReg struct {
	registerslib.Byte

	Shades []*DMGPaletteColorShade
}

func (r *DMGPaletteReg) GetColor(id uint8) frontends.Pixel {
	paletteShade := r.Shades[id]
	if paletteShade == nil {
		return frontends.TransparentPixel
	}
	return paletteShade.GetColor()
}
