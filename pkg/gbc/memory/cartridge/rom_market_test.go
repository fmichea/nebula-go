package cartridge

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"nebula-go/pkg/gbc/memory/lib"
)

func TestROMMarket(t *testing.T) {
	setup := func(val uint8) []uint8 {
		romData := make([]uint8, _minimumROMDataSize)
		romData[_targetMarketAddress] = val
		return romData
	}

	t.Run("ROM for japanese market is detected", func(t *testing.T) {
		assert.Equal(t, lib.Japanese, loadROMMarket(setup(0x00)))
	})

	t.Run("ROM for non-japanese market is detected", func(t *testing.T) {
		assert.Equal(t, lib.NonJapanese, loadROMMarket(setup(0x01)))
	})
}
