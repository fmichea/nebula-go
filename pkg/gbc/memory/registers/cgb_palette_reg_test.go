package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"nebula-go/pkg/common/frontends"
)

func TestCGBPaletteReg(t *testing.T) {
	setup := func() (*CGBPaletteIndexReg, *CGBPaletteReg) {
		var value uint8
		index := NewCGBPaletteIndexReg(&value)
		return index, NewCGBPaletteReg(index)
	}

	t.Run("default color is white", func(t *testing.T) {
		index, reg := setup()

		assert.Equal(t, uint8(0x00), index.Get())
		assert.Equal(t, uint8(0xFF), reg.Get())

		index.Index.Set(0x01)
		assert.Equal(t, uint8(0x7F), reg.Get())

		assert.Equal(t, frontends.White, reg.GetColor(0, 0))
	})

	t.Run("set uses index to select which byte to change: low byte affects red", func(t *testing.T) {
		index, reg := setup()

		// palette 1, color 1, byte 0
		index.Index.Set(10)
		reg.Set(0xE0)

		assert.Equal(t, frontends.White, reg.GetColor(0, 0))
		assert.Equal(t, frontends.NewPixel(0x00, 0xFF, 0xFF), reg.GetColor(1, 1))
	})

	t.Run("set uses index to select which byte to change: high byte affects blue", func(t *testing.T) {
		index, reg := setup()

		// palette 4, color 0, byte 1
		index.Index.Set(33)
		reg.Set(0x03)

		assert.Equal(t, frontends.White, reg.GetColor(0, 0))
		assert.Equal(t, frontends.NewPixel(0xFF, 0xFF, 0x00), reg.GetColor(4, 0))
	})

	t.Run("set uses index to select which byte to change: high and low affect green", func(t *testing.T) {
		index, reg := setup()

		// palette 6, color 3, byte 0
		index.Index.Set(54)
		reg.Set(0x1F)

		// palette 6, color 3, byte 1
		index.Index.Set(55)
		reg.Set(0x7C)

		assert.Equal(t, frontends.White, reg.GetColor(0, 0))
		assert.Equal(t, frontends.NewPixel(0xFF, 0x00, 0xFF), reg.GetColor(6, 3))
	})

	t.Run("set with auto-increment changes index: normal case", func(t *testing.T) {
		index, reg := setup()

		assert.Equal(t, frontends.White, reg.GetColor(0, 0))
		assert.Equal(t, frontends.White, reg.GetColor(0, 1))

		index.AutoIncrement.SetBool(true)
		reg.Set(0x00)
		reg.Set(0x00)

		assert.Equal(t, uint8(0x82), index.Get())
		assert.Equal(t, frontends.Black, reg.GetColor(0, 0))
		assert.Equal(t, frontends.White, reg.GetColor(0, 1))
	})

	t.Run("set with auto-increment changes index: last index", func(t *testing.T) {
		index, reg := setup()

		assert.Equal(t, frontends.White, reg.GetColor(7, 3))
		assert.Equal(t, frontends.White, reg.GetColor(0, 0))

		index.AutoIncrement.SetBool(true)
		index.Index.Set(0x3E)
		reg.Set(0x00)
		reg.Set(0x00)

		assert.Equal(t, uint8(0x80), index.Get())
		assert.Equal(t, frontends.Black, reg.GetColor(7, 3))
		assert.Equal(t, frontends.White, reg.GetColor(0, 0))
	})

	t.Run("color shade is adjusted", func(t *testing.T) {
		_, reg := setup()

		reg.Set(0xEC)

		assert.Equal(t, frontends.NewPixel(0x62, 0xFF, 0xFF), reg.GetColor(0, 0))
	})
}
