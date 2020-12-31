package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

func TestDMGPaletteColorShade(t *testing.T) {
	setup := func() (registerslib.Byte, *DMGPaletteColorShade) {
		var value uint8

		reg := registerslib.NewByte(&value, 0x00)
		return reg, NewDMGPaletteColorShade(reg, 0x0)
	}

	t.Run("shade 0 is white", func(t *testing.T) {
		_, shade := setup()

		assert.Equal(t, uint8(0), shade.Get())
		assert.Equal(t, _white, shade.GetColor())
	})

	t.Run("value outside of shade's offset is ignored", func(t *testing.T) {
		reg, shade := setup()

		reg.Set(0xFC)
		assert.Equal(t, uint8(0), shade.Get())
		assert.Equal(t, _white, shade.GetColor())
	})

	t.Run("shade 1 is light grey", func(t *testing.T) {
		reg, shade := setup()

		reg.Set(0x01)
		assert.Equal(t, uint8(0x1), shade.Get())
		assert.Equal(t, _lightGrey, shade.GetColor())
	})

	t.Run("shade 2 is dark grey", func(t *testing.T) {
		reg, shade := setup()

		reg.Set(0x02)
		assert.Equal(t, uint8(0x2), shade.Get())
		assert.Equal(t, _darkGrey, shade.GetColor())
	})

	t.Run("shade 3 is dark grey", func(t *testing.T) {
		reg, shade := setup()

		reg.Set(0x03)
		assert.Equal(t, uint8(0x3), shade.Get())
		assert.Equal(t, _black, shade.GetColor())
	})
}
