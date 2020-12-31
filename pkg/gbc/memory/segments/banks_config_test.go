package segments

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A lot of this is also black box tested through the general segment tests.
func TestBanksConfig(t *testing.T) {
	t.Run("default is one bank, bank 0 selected, unpinned", func(t *testing.T) {
		bc := newBanksConfig(0x100, 0x1FF)
		require.NoError(t, bc.validateAndInitialize())

		t.Run("basic configuration", func(t *testing.T) {
			assert.Equal(t, bc.segmentAddressRange, bc.bankedAddressRange)
			assert.Equal(t, uint(0), bc.current)
			assert.Equal(t, uint(1), bc.count)
			assert.False(t, bc.bank0Pinned)
		})

		t.Run("banks run the whole segment", func(t *testing.T) {
			assert.Equal(t, uint(0x100), bc.sizePerBank())
		})

		t.Run("contains address on the whole segment", func(t *testing.T) {
			assert.True(t, bc.containsAddress(0x100))
			assert.True(t, bc.containsAddress(0x180))
			assert.True(t, bc.containsAddress(0x1FF))
			assert.False(t, bc.containsAddress(0x200))
		})

		t.Run("asOffset uses the beginning of whole segment", func(t *testing.T) {
			assert.Equal(t, uint(0x10), bc.asOffset(0x110))
			assert.Equal(t, uint(0x18), bc.asOffset(0x118))
		})

		t.Run("cannot select bank outside of bounds", func(t *testing.T) {
			assert.Equal(t, ErrBankUnavailable, bc.selectBank(123))
		})

		t.Run("is address banked", func(t *testing.T) {
			assert.False(t, bc.isBanked(0x100))
			assert.False(t, bc.isBanked(0x1FF))
		})
	})

	t.Run("cannot pin bank 0 on only one bank", func(t *testing.T) {
		bc := newBanksConfig(0x100, 0x1FF)
		bc.makeBank0Pinned()

		assert.Equal(t, ErrCannotPin0WithOneBank, bc.validateAndInitialize())
	})

	t.Run("cannot set bank count to 0", func(t *testing.T) {
		bc := newBanksConfig(0x100, 0x1FF)
		assert.Equal(t, ErrBankCountInvalid, bc.setBankCount(0))
	})

	t.Run("with a few banks, no pinning", func(t *testing.T) {
		bc := newBanksConfig(0x100, 0x1FF)
		require.NoError(t, bc.setBankCount(5))
		require.NoError(t, bc.validateAndInitialize())

		t.Run("basic configuration", func(t *testing.T) {
			assert.Equal(t, bc.segmentAddressRange, bc.bankedAddressRange)
			assert.Equal(t, uint(0), bc.current)
			assert.Equal(t, uint(5), bc.count)
			assert.False(t, bc.bank0Pinned)
		})

		t.Run("contains address on the whole segment", func(t *testing.T) {
			assert.True(t, bc.containsAddress(0x100))
			assert.True(t, bc.containsAddress(0x180))
			assert.True(t, bc.containsAddress(0x1FF))
			assert.False(t, bc.containsAddress(0x200))
		})

		t.Run("select bank works", func(t *testing.T) {
			assert.NoError(t, bc.selectBank(2))
			assert.Equal(t, uint(2), bc.current)
		})

		t.Run("is address banked", func(t *testing.T) {
			assert.True(t, bc.isBanked(0x100))
			assert.True(t, bc.isBanked(0x1FF))
		})
	})

	t.Run("bank 0 pinned halves the banking zone", func(t *testing.T) {
		bc := newBanksConfig(0x100, 0x1FF)
		require.NoError(t, bc.setBankCount(5))
		bc.makeBank0Pinned()

		require.NoError(t, bc.validateAndInitialize())

		t.Run("basic configuration", func(t *testing.T) {
			assert.NotEqual(t, bc.segmentAddressRange, bc.bankedAddressRange)
			assert.Equal(t, uint(1), bc.current)
			assert.Equal(t, uint(5), bc.count)
			assert.True(t, bc.bank0Pinned)
		})

		t.Run("size per bank is now half of full segment size", func(t *testing.T) {
			assert.Equal(t, uint(0x80), bc.sizePerBank())
		})

		t.Run("only contains address in second half", func(t *testing.T) {
			assert.False(t, bc.containsAddress(0x100))
			assert.False(t, bc.containsAddress(0x17F))
			assert.True(t, bc.containsAddress(0x180))
			assert.True(t, bc.containsAddress(0x1FF))
			assert.False(t, bc.containsAddress(0x200))
		})

		t.Run("offset starts at second half", func(t *testing.T) {
			assert.Equal(t, uint(0x10), bc.asOffset(0x190))
			assert.Equal(t, uint(0x20), bc.asOffset(0x1A0))
		})

		t.Run("selecting bank 0 actually selects 1", func(t *testing.T) {
			assert.NoError(t, bc.selectBank(0))
			assert.Equal(t, uint(1), bc.current)
		})

		t.Run("is address banked", func(t *testing.T) {
			assert.False(t, bc.isBanked(0x100))
			assert.True(t, bc.isBanked(0x1FF))
		})
	})
}
