package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	_invalidROMType = ROMType(0xFF)
)

func TestROMType_String(t *testing.T) {
	t.Run("DMG has string value", func(t *testing.T) {
		assert.Equal(t, "GAME BOY (DMG-01)", DMG01.String())
	})

	t.Run("GBC has string value", func(t *testing.T) {
		assert.Equal(t, "GAME BOY COLOR (CGB-001)", CGB001.String())
	})

	t.Run("unknown value has string", func(t *testing.T) {
		assert.Equal(t, "UNKNOWN", _invalidROMType.String())
	})
}

func TestROMType_GetTitleSize(t *testing.T) {
	t.Run("DMG has 0x16 characters long cartridge title", func(t *testing.T) {
		assert.Equal(t, 0x16, DMG01.GetTitleSize())
	})

	t.Run("CGB has reduced 0x15 characters long cartridge title", func(t *testing.T) {
		assert.Equal(t, 0x15, CGB001.GetTitleSize())
	})
}

func TestROMType_VRAMBankCount(t *testing.T) {
	t.Run("DMG does not have VRAM banking", func(t *testing.T) {
		assert.Equal(t, uint(1), DMG01.VRAMBankCount())
	})

	t.Run("CGB has 2 banks for VRAM", func(t *testing.T) {
		assert.Equal(t, uint(2), CGB001.VRAMBankCount())
	})
}

func TestROMType_WRAMBankCount(t *testing.T) {
	t.Run("DMG does not have WRAM banking (bank0 is pinned)", func(t *testing.T) {
		assert.Equal(t, uint(2), DMG01.WRAMBankCount())
	})

	t.Run("CGB has 8 banks available for WRAM (bank0 is pinned)", func(t *testing.T) {
		assert.Equal(t, uint(8), CGB001.WRAMBankCount())
	})
}
