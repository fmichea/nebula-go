package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterruptReg(t *testing.T) {
	t.Run("all is disabled by default", func(t *testing.T) {
		var value uint8

		reg := NewInterruptReg(&value)
		assert.False(t, reg.VBlank.IsRequested())
		assert.False(t, reg.STAT.IsRequested())
		assert.False(t, reg.Timer.IsRequested())
		assert.False(t, reg.Serial.IsRequested())
		assert.False(t, reg.Joypad.IsRequested())
	})

	t.Run("some of them can be disabled", func(t *testing.T) {
		var value uint8

		reg := NewInterruptReg(&value)
		reg.Set(0b10101)

		assert.True(t, reg.VBlank.IsRequested())
		assert.False(t, reg.STAT.IsRequested())
		assert.True(t, reg.Timer.IsRequested())
		assert.False(t, reg.Serial.IsRequested())
		assert.True(t, reg.Joypad.IsRequested())
	})

	t.Run("interrupt flag can be requested", func(t *testing.T) {
		var value uint8

		reg := NewInterruptReg(&value)
		assert.Equal(t, uint8(0x00), reg.Get())

		reg.STAT.Request()
		assert.Equal(t, uint8(0x02), reg.Get())
	})

	t.Run("interrupt flag can be acknowledged", func(t *testing.T) {
		var value uint8

		reg := NewInterruptReg(&value)
		reg.Set(0x0F)

		assert.Equal(t, uint8(0x0F), reg.Get())

		reg.STAT.Acknowledge()
		assert.Equal(t, uint8(0x0D), reg.Get())
	})
}
