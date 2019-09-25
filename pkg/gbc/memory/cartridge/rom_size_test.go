package cartridge

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"nebula-go/pkg/gbc/memory/lib"
)

func TestROMSize(t *testing.T) {
	setup := func(val uint8) []uint8 {
		romData := make([]uint8, _minimumROMDataSize)
		romData[0x148] = val
		return romData
	}

	t.Run("invalid value does not pass check", func(t *testing.T) {
		assert.Equal(t, ErrROMSizeInvalid, verifyROMSizeFlag(setup(0xFF)))
	})

	t.Run("valid value passes check and can be loaded", func(t *testing.T) {
		romData := setup(0x02)
		assert.NoError(t, verifyROMSizeFlag(romData))

		romSize := loadROMSize(romData)
		assert.True(t, romSize.IsValid())
		assert.Equal(t, lib.ROMSize128KB, romSize)
	})
}
