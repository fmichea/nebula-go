package bitwise

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInverseMask8(t *testing.T) {
	t.Run("nothing masked", func(t *testing.T) {
		assert.Equal(t, uint8(0xFF), InverseMask8(0xFF, 0x00))
	})

	t.Run("mask the bottom nibble", func(t *testing.T) {
		assert.Equal(t, uint8(0xF0), InverseMask8(0xFF, 0x0F))
	})

	t.Run("mask the top nibble", func(t *testing.T) {
		assert.Equal(t, uint8(0x0F), InverseMask8(0xFF, 0xF0))
	})

	t.Run("advanced mask", func(t *testing.T) {
		assert.Equal(t, uint8(0x50), InverseMask8(0xD2, 0xA2))
	})
}

func TestInverseMask16(t *testing.T) {
	t.Run("nothing masked", func(t *testing.T) {
		assert.Equal(t, uint16(0xFFFF), InverseMask16(0xFFFF, 0x0000))
	})

	t.Run("mask the first nibble from bottom", func(t *testing.T) {
		assert.Equal(t, uint16(0xFFF0), InverseMask16(0xFFFF, 0x000F))
	})

	t.Run("mask the second nibble from bottom", func(t *testing.T) {
		assert.Equal(t, uint16(0xFF0F), InverseMask16(0xFFFF, 0x00F0))
	})
}

func TestInverseMask32(t *testing.T) {
	t.Run("nothing masked", func(t *testing.T) {
		assert.Equal(t, uint32(0xFFFFFFFF), InverseMask32(0xFFFFFFFF, 0x00000000))
	})

	t.Run("mask the first nibble from bottom", func(t *testing.T) {
		assert.Equal(t, uint32(0xFFFFFFF0), InverseMask32(0xFFFFFFFF, 0x0000000F))
	})

	t.Run("mask the second nibble from bottom", func(t *testing.T) {
		assert.Equal(t, uint32(0xFFFFFF0F), InverseMask32(0xFFFFFFFF, 0x000000F0))
	})
}
