package cartridge

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"nebula-go/pkg/gbc/memory/lib"
)

func TestROMTitle(t *testing.T) {
	longTitle := "GAME TITLE IS TOO LONG"
	smallTitle := "GAME"

	setup := func(str string) []uint8 {
		romData := make([]uint8, _romSizeAddress)
		for idx, value := range []uint8(str) {
			offset := _titleAddress + idx
			if len(romData) <= offset {
				continue
			}
			romData[offset] = value
		}
		return romData
	}

	t.Run("DMG has a title of 16 characters long", func(t *testing.T) {
		assert.Equal(t, "GAME TITLE IS TO", loadROMTitle(setup(longTitle), lib.DMG01))
		assert.Equal(t, "GAME", loadROMTitle(setup(smallTitle), lib.DMG01))
	})

	t.Run("CGB has a title of 15 characters long", func(t *testing.T) {
		assert.Equal(t, "GAME TITLE IS T", loadROMTitle(setup(longTitle), lib.CGB001))
		assert.Equal(t, "GAME", loadROMTitle(setup(smallTitle), lib.CGB001))
	})
}
