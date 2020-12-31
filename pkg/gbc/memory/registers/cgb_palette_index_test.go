package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCGBPaletteIndexReg(t *testing.T) {
	t.Run("default value for reg is 0", func(t *testing.T) {
		value := uint8(0xFF)
		reg := NewCGBPaletteIndexReg(&value)

		assert.Equal(t, uint8(0x00), reg.Get())
		assert.Equal(t, uint8(0x00), value)
	})

	t.Run("setting auto-increment sets bit 7", func(t *testing.T) {
		var value uint8

		reg := NewCGBPaletteIndexReg(&value)
		assert.Equal(t, uint8(0x00), reg.Get())
		assert.False(t, reg.AutoIncrement.GetBool())

		reg.AutoIncrement.SetBool(true)
		assert.Equal(t, uint8(0x80), reg.Get())
		assert.True(t, reg.AutoIncrement.GetBool())
	})

	t.Run("index covers lower 5 bits", func(t *testing.T) {
		var value uint8

		reg := NewCGBPaletteIndexReg(&value)
		assert.Equal(t, uint8(0x00), reg.Get())
		assert.Equal(t, uint8(0x00), reg.Index.Get())

		reg.Index.Set(0xFF)
		assert.Equal(t, uint8(0x3F), reg.Index.Get())
	})
}
