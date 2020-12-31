package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKEY1Reg(t *testing.T) {
	uint8ptr := func() *uint8 {
		var value uint8
		return &value
	}

	t.Run("default speed is normal", func(t *testing.T) {
		reg := NewKEY1Reg(uint8ptr())
		assert.Equal(t, uint8(0x00), reg.Get())
		assert.False(t, reg.CurrentSpeed.IsDoubleSpeed())
	})

	t.Run("no mode switch request has no effect", func(t *testing.T) {
		reg := NewKEY1Reg(uint8ptr())
		assert.Equal(t, uint8(0x00), reg.Get())

		reg.SwitchIfRequested()
		assert.Equal(t, uint8(0x00), reg.Get())
	})

	t.Run("switch mode request toggles the speed", func(t *testing.T) {
		reg := NewKEY1Reg(uint8ptr())
		reg.Set(0x01)
		assert.Equal(t, uint8(0x01), reg.Get())

		reg.SwitchIfRequested()
		assert.Equal(t, uint8(0x80), reg.Get())
		assert.True(t, reg.CurrentSpeed.IsDoubleSpeed())

		reg.Set(0x01)
		assert.Equal(t, uint8(0x81), reg.Get())

		reg.SwitchIfRequested()
		assert.Equal(t, uint8(0x00), reg.Get())
		assert.False(t, reg.CurrentSpeed.IsDoubleSpeed())
	})

	t.Run("KEY1 reg's 7-th bit can only be changed using the flag", func(t *testing.T) {
		reg := NewKEY1Reg(uint8ptr())
		reg.Set(0xFF)

		assert.Equal(t, uint8(0x01), reg.Get())

		reg.CurrentSpeed.SetBool(true)
		assert.Equal(t, uint8(0x81), reg.Get())
	})
}
