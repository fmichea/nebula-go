package graphics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTileAttributes(t *testing.T) {
	attrs := NewTileAttributes(0b10101010)

	require.NotNil(t, attrs)
	assert.Equal(t, uint8(0b010), attrs.CGBPalette.Get())
	assert.Equal(t, uint8(0b1), attrs.VRAMBank.Get())
	assert.Equal(t, uint8(0b0), attrs.DMGPalette.Get())
	assert.True(t, attrs.HFlip.GetBool())
	assert.False(t, attrs.VFlip.GetBool())
	assert.True(t, attrs.BackgroundPriority.GetBool())
}
