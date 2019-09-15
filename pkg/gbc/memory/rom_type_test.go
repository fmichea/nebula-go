package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestROMType_String(t *testing.T) {
	t.Run("DMG game-boy rom type has string value", func(t *testing.T) {
		assert.Equal(t, "GAME BOY (DMG-01)", DMG01.String())
	})

	t.Run("Game Boy Color rom type has string value", func(t *testing.T) {
		assert.Equal(t, "GAME BOY COLOR (CGB-001)", CGB001.String())
	})

	t.Run("unknown value has string", func(t *testing.T) {
		assert.Equal(t, "UNKNOWN", ROMType(0xff).String())
	})
}
