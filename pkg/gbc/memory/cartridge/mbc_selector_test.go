package cartridge

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"nebula-go/pkg/gbc/memory/lib"
)

func TestMBCSelector(t *testing.T) {
	setup := func(val uint8) []uint8 {
		romData := make([]uint8, _minimumROMDataSize)
		romData[0x147] = val
		return romData
	}

	t.Run("invalid MBC selector is denied", func(t *testing.T) {
		romData := setup(040)
		assert.Equal(t, lib.ErrMBCNotImplemented, verifyMBCSelector(romData))
	})

	t.Run("valid MBC selector can be verified and loaded", func(t *testing.T) {
		romData := setup(0x00)
		assert.NoError(t, verifyMBCSelector(romData))

		selector := loadMBCSelector(romData)
		assert.NotNil(t, selector)
		assert.True(t, selector.IsValid())
	})
}
