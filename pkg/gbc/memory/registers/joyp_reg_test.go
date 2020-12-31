package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJOYPReg(t *testing.T) {
	t.Run("default value is FFh", func(t *testing.T) {
		joyp := NewJOYPReg(nil)
		assert.Equal(t, uint8(0xFF), joyp.Get())
	})

	t.Run("selecting direction keys using bit4=0, no button selected", func(t *testing.T) {
		joyp := NewJOYPReg(nil)
		joyp.Set(0x20)

		assert.Equal(t, uint8(0xEF), joyp.Get())
	})

	t.Run("selecting button keys using bit5=0, no button selected", func(t *testing.T) {
		joyp := NewJOYPReg(nil)
		joyp.Set(0x10)

		assert.Equal(t, uint8(0xDF), joyp.Get())
	})

	t.Run("selecting direction keys, down+a button pressed", func(t *testing.T) {
		joyp := NewJOYPReg(nil)
		joyp.DownButton.Press()
		joyp.AButton.Press()

		joyp.Set(0x20)

		assert.Equal(t, uint8(0xE7), joyp.Get())
	})

	t.Run("selecting button keys, down+a button pressed", func(t *testing.T) {
		joyp := NewJOYPReg(nil)
		joyp.DownButton.Press()
		joyp.AButton.Press()

		joyp.Set(0x10)

		assert.Equal(t, uint8(0xDE), joyp.Get())
	})

	t.Run("selecting direction keys and button keys, directions returned", func(t *testing.T) {
		joyp := NewJOYPReg(nil)
		joyp.DownButton.Press()
		joyp.AButton.Press()

		joyp.Set(0x00)

		assert.Equal(t, uint8(0xC7), joyp.Get())
	})

	t.Run("button can be released", func(t *testing.T) {
		joyp := NewJOYPReg(nil)
		joyp.Set(0x20)

		joyp.DownButton.Press()
		assert.Equal(t, uint8(0xE7), joyp.Get())

		joyp.DownButton.Release()
		assert.Equal(t, uint8(0xEF), joyp.Get())
	})
}
