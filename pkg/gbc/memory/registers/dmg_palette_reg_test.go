package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"nebula-go/pkg/common/frontends"
)

func TestDMGPaletteReg(t *testing.T) {
	setup := func(value uint8, transparent0 bool) *DMGPaletteReg {
		return NewDMGPaletteReg(&value, value, transparent0)
	}

	t.Run("colorID's shade is used to get the color", func(t *testing.T) {
		reg := setup(0b00011011, false)

		assert.Equal(t, _black, reg.GetColor(0))
		assert.Equal(t, _darkGrey, reg.GetColor(1))
		assert.Equal(t, _lightGrey, reg.GetColor(2))
		assert.Equal(t, _white, reg.GetColor(3))
	})

	t.Run("colodID 0 is transparent when transparent0 is set", func(t *testing.T) {
		reg := setup(0x00, true)

		assert.Equal(t, frontends.TransparentPixel, reg.GetColor(0))
		assert.Equal(t, _white, reg.GetColor(1))
		assert.Equal(t, _white, reg.GetColor(2))
		assert.Equal(t, _white, reg.GetColor(3))
	})
}
