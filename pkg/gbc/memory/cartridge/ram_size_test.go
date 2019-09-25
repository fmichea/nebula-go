package cartridge

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"nebula-go/pkg/gbc/memory/lib"
)

func TestRAMSize(t *testing.T) {
	setup := func(val uint8) []uint8 {
		romData := make([]uint8, _minimumROMDataSize)
		romData[0x149] = val
		return romData
	}

	t.Run("invalid value does not pass check", func(t *testing.T) {
		assert.Equal(t, ErrRAMSizeInvalid, verifyRAMSizeFlag(setup(0xFF)))
	})

	t.Run("valid value passes check and can be loaded", func(t *testing.T) {
		romData := setup(0x00)
		assert.NoError(t, verifyRAMSizeFlag(romData))

		ramSize := loadRAMSize(romData)
		assert.True(t, ramSize.IsValid())
		assert.Equal(t, lib.RAMSizeNone, ramSize)
	})
}
