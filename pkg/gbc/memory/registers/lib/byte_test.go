package registerslib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteReg_Get(t *testing.T) {
	t.Run("init changes pointer value", func(t *testing.T) {
		var value uint8

		assert.NotEqual(t, uint8(0xFF), value)
		_ = NewByte(&value, 0xFF)
		assert.Equal(t, uint8(0xFF), value)
	})

	t.Run("get returns the current pointer's value", func(t *testing.T) {
		var value uint8

		reg := NewByte(&value, 0xFF)

		value = 0x55
		assert.Equal(t, uint8(0x55), reg.Get())
	})

	t.Run("set changes value", func(t *testing.T) {
		var value uint8

		reg := NewByte(&value, 0xFF)
		reg.Set(0x55)

		assert.Equal(t, uint8(0x55), value)
	})

	t.Run("mask is used only when using set, not original value", func(t *testing.T) {
		var value uint8

		reg := NewByteWithMask(&value, 0x45, 0x0F)
		assert.Equal(t, uint8(0x45), reg.Get())

		reg.Set(0xF8)
		assert.Equal(t, uint8(0x48), reg.Get())
	})
}
