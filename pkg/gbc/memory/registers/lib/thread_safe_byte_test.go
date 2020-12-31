package registerslib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreadSafeByte(t *testing.T) {
	t.Run("set changes value", func(t *testing.T) {
		reg := NewThreadSafeByte(0xAB)
		assert.Equal(t, uint8(0xAB), reg.Get())

		reg.Set(0xCD)
		assert.Equal(t, uint8(0xCD), reg.Get())
	})

	t.Run("with mask can be used", func(t *testing.T) {
		reg := NewThreadSafeByteWithMask(0xAB, 0x0F)
		assert.Equal(t, uint8(0xAB), reg.Get())

		reg.Set(0xF8)
		assert.Equal(t, uint8(0xA8), reg.Get())
	})
}
