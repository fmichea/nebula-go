package cartridge

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"nebula-go/pkg/gbc/memory/lib"
)

func TestROMType(t *testing.T) {
	setup := func(val uint8) []uint8 {
		romData := make([]uint8, _minimumROMDataSize)
		romData[0x143] = val
		return romData
	}

	t.Run("CGB values are loaded properly", func(t *testing.T) {
		assert.Equal(t, lib.CGB001, loadROMType(setup(0x80)))
		assert.Equal(t, lib.CGB001, loadROMType(setup(0xC0)))
	})

	t.Run("DMG values are loaded properly", func(t *testing.T) {
		assert.Equal(t, lib.DMG01, loadROMType(setup(0)))
		assert.Equal(t, lib.DMG01, loadROMType(setup(0x45)))
	})
}
