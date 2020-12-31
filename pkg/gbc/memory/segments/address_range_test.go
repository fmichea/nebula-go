package segments

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressRange(t *testing.T) {
	ar := AddressRange{Start: 0x00, End: 0xFF}
	ar2 := AddressRange{Start: 0x100, End: 0x1FF}

	t.Run("size is number of bytes between start and end", func(t *testing.T) {
		assert.Equal(t, uint(0x100), ar.size())
	})

	t.Run("contains address works properly", func(t *testing.T) {
		assert.True(t, ar.containsAddress(0))
		assert.True(t, ar.containsAddress(0x80))
		assert.True(t, ar.containsAddress(0xFF))
		assert.False(t, ar.containsAddress(0x100))
	})

	t.Run("transposeAddress moves a given address in original AddressRange to other one", func(t *testing.T) {
		assert.Equal(t, uint16(0x110), ar.transposeAddress(ar2, 0x10))
	})

	t.Run("asOffset returns address as if it was 0 indexed", func(t *testing.T) {
		assert.Equal(t, uint(0x10), ar.asOffset(0x10))
		assert.Equal(t, uint(0x10), ar2.asOffset(0x110))
	})
}
