package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSTATTimer(t *testing.T) {
	t.Run("forward decrements the wait time", func(t *testing.T) {
		timer := NewSTATTimer()

		assert.Equal(t, int16(0), timer.waitCount)

		timer.Forward(4)
		assert.Equal(t, int16(-4), timer.waitCount)

		timer.Forward(8)
		assert.Equal(t, int16(-12), timer.waitCount)
	})

	t.Run("Expired is set when wait is done", func(t *testing.T) {
		timer := NewSTATTimer()

		assert.Equal(t, int16(0), timer.waitCount)
		assert.False(t, timer.Expired())

		timer.Forward(1)
		assert.True(t, timer.Expired())
	})

	t.Run("changing mode adds wait time", func(t *testing.T) {
		timer := NewSTATTimer()
		timer.Forward(4)
		timer.SwitchMode(STATModeVBlank)

		assert.False(t, timer.Expired())
	})
}
