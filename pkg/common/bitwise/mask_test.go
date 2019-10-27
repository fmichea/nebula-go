package bitwise

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMask8(t *testing.T) {
	t.Run("nothing masked", func(t *testing.T) {
		assert.Equal(t, uint8(0xFF), Mask8(0xFF, 0xFF))
	})

	t.Run("first nibble masked", func(t *testing.T) {
		assert.Equal(t, uint8(0xF0), Mask8(0xFF, 0xF0))
	})

	t.Run("second nibble masked", func(t *testing.T) {
		assert.Equal(t, uint8(0x0F), Mask8(0xFF, 0x0F))
	})

	t.Run("advanced mask", func(t *testing.T) {
		assert.Equal(t, uint8(0x50), Mask8(0x52, 0xF0))
	})
}

func TestMask16(t *testing.T) {
	t.Run("nothing masked", func(t *testing.T) {
		assert.Equal(t, uint16(0xFFFF), Mask16(0xFFFF, 0xFFFF))
	})

	t.Run("first nibble masked", func(t *testing.T) {
		assert.Equal(t, uint16(0xFFF0), Mask16(0xFFFF, 0xFFF0))
	})

	t.Run("middle nibbles masked", func(t *testing.T) {
		assert.Equal(t, uint16(0xF00F), Mask16(0xFFFF, 0xF00F))
	})
}

func TestMask32(t *testing.T) {
	t.Run("nothing masked", func(t *testing.T) {
		assert.Equal(t, uint32(0xFFFF), Mask32(0xFFFF, 0xFFFF))
	})

	t.Run("first nibble masked", func(t *testing.T) {
		assert.Equal(t, uint32(0xFFF0), Mask32(0xFFFF, 0xFFF0))
	})

	t.Run("middle nibbles of first 16 bits masked", func(t *testing.T) {
		assert.Equal(t, uint32(0xF00F), Mask32(0xFFFF, 0xF00F))
	})
}
