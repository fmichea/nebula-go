package bitfields

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByte_Get(t *testing.T) {
	t.Run("simple byte has initial value", func(t *testing.T) {
		reg := NewByte(0xAB)
		assert.Equal(t, uint8(0xAB), reg.Get())
	})

	t.Run("byte with mask has initial value masked", func(t *testing.T) {
		reg := NewByteWithMask(0xAB, 0x0F)
		assert.Equal(t, uint8(0x0B), reg.Get())
	})
}

func TestByte_Set(t *testing.T) {
	t.Run("simple byte can have value changed", func(t *testing.T) {
		reg := NewByte(0xAB)
		reg.Set(0xDE)

		assert.Equal(t, uint8(0xDE), reg.Get())
	})

	t.Run("byte with mask has initial value masked", func(t *testing.T) {
		reg := NewByteWithMask(0xAB, 0x0F)
		reg.Set(0xDE)

		assert.Equal(t, uint8(0x0E), reg.Get())
	})
}
