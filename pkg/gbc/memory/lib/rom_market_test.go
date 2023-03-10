package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestROMMarket_String(t *testing.T) {
	t.Run("japanese ROM is detected", func(t *testing.T) {
		value := ROMMarket(0x00)
		assert.Equal(t, "Japanese", value.String())
	})

	t.Run("non-japanese ROM is detected", func(t *testing.T) {
		value := ROMMarket(0x01)
		assert.Equal(t, "Non-Japanese", value.String())
	})

	t.Run("undefined value is detected", func(t *testing.T) {
		value := ROMMarket(0x13)
		assert.Equal(t, "Unknown", value.String())
	})
}
