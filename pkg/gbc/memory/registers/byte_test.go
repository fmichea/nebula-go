package registers

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

	t.Run("mapped byte works the same", func(t *testing.T) {
		b := uint8(0xAB)
		reg := NewMappedByte(&b)

		assert.Equal(t, b, reg.Get())
	})

	t.Run("mapped byte with mask is masked at init", func(t *testing.T) {
		v := uint8(0xAB)
		reg := NewMappedByteWithMask(&v, 0xF0)

		expected := uint8(0xA0)
		assert.Equal(t, expected, reg.Get())
		assert.Equal(t, expected, v)
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

	t.Run("mapped byte works the same", func(t *testing.T) {
		b := uint8(0xAB)

		expected := uint8(0xDE)

		reg := NewMappedByte(&b)
		reg.Set(expected)

		assert.Equal(t, expected, reg.Get())
		assert.Equal(t, expected, b)
	})

	t.Run("mapped byte with mask is masked at init", func(t *testing.T) {
		b := uint8(0xAB)

		reg := NewMappedByteWithMask(&b, 0xF0)
		reg.Set(0xDE)

		expected := uint8(0xD0)
		assert.Equal(t, expected, reg.Get())
		assert.Equal(t, expected, b)
	})
}
