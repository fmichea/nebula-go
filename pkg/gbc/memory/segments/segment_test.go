package segments

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testify's assert.Equal and assert.NotEqual do not work as expected for pointers, as they test the pointed value's
// (non-)equality, instead of the pointers themselves.
func assertPtrEqual(t *testing.T, ptr1, ptr2 *uint8) {
	assert.True(t, ptr1 == ptr2, "Pointers %p and %p are not equal", ptr1, ptr2)
}

func assertPtrNotEqual(t *testing.T, ptr1, ptr2 *uint8) {
	assert.True(t, ptr1 != ptr2, "Pointers are equal with value %p", ptr1)
}

func TestNew(t *testing.T) {
	setup := func(t *testing.T) Segment {
		s, err := New(0x00, 0xFF)
		require.NoError(t, err)
		require.NotNil(t, s)
		return s
	}

	t.Run("contains address", func(t *testing.T) {
		s := setup(t)

		assert.True(t, s.ContainsAddress(0x00))
		assert.True(t, s.ContainsAddress(0x80))
		assert.True(t, s.ContainsAddress(0xFF))
		assert.False(t, s.ContainsAddress(0x100))
	})

	t.Run("ptr is the same for a given address", func(t *testing.T) {
		s := setup(t)

		ptr := s.BytePtr(0x01)
		require.NotNil(t, ptr)
		assert.Equal(t, uint8(0), *ptr)

		*ptr = 0xff
		assert.Equal(t, uint8(0xff), *s.BytePtr(0x01))
	})

	t.Run("can only select bank 0", func(t *testing.T) {
		s := setup(t)

		assert.Equal(t, ErrBankUnavailable, s.SelectBank(1))
		assert.Equal(t, ErrBankUnavailable, s.SelectBank(13))
		assert.NoError(t, s.SelectBank(0))
	})

	t.Run("pointer outside range is nil", func(t *testing.T) {
		s := setup(t)
		assert.Nil(t, s.BytePtr(0x1000))
	})
}

func TestNewWithBanks(t *testing.T) {
	setup := func(t *testing.T) Segment {
		s, err := New(0, 0xFF, WithBanks(3))
		require.NoError(t, err)
		require.NotNil(t, s)
		return s
	}

	t.Run("contains address works the same as normal segment", func(t *testing.T) {
		s := setup(t)

		assert.True(t, s.ContainsAddress(0x00))
		assert.True(t, s.ContainsAddress(0x80))
		assert.True(t, s.ContainsAddress(0xFF))
		assert.False(t, s.ContainsAddress(0x100))
	})

	t.Run("ptr is the same for a given address without bank change", func(t *testing.T) {
		s := setup(t)

		ptr := s.BytePtr(0x01)
		assert.Equal(t, uint8(0), *ptr)

		*ptr = 0xff
		assert.Equal(t, uint8(0xff), *s.BytePtr(0x01))
	})

	t.Run("bank 0 is the default", func(t *testing.T) {
		s := setup(t)

		ptr := s.BytePtr(0x01)
		require.NotNil(t, ptr)

		assert.NoError(t, s.SelectBank(0))
		ptr2 := s.BytePtr(0x01)

		assertPtrEqual(t, ptr, ptr2)
	})

	t.Run("bank changes change the actual memory values", func(t *testing.T) {
		s := setup(t)

		// by default we should be on bank 0.
		ptr := s.BytePtr(0x01)
		require.NotNil(t, ptr)
		assert.Equal(t, uint8(0), *ptr)

		*ptr = 0xff

		// moving to bank 1 and checking we did change memory view.
		assert.NoError(t, s.SelectBank(1))

		ptr2 := s.BytePtr(0x01)
		require.NotNil(t, ptr2)
		assertPtrNotEqual(t, ptr, ptr2)
		assert.Equal(t, uint8(0), *ptr2)

		// switching back to bank 0.
		assert.NoError(t, s.SelectBank(0))

		ptr3 := s.BytePtr(0x01)
		require.NotNil(t, ptr3)
		assertPtrNotEqual(t, ptr2, ptr3)
		assertPtrEqual(t, ptr3, ptr)
		assert.Equal(t, uint8(0xff), *ptr3)
	})

	t.Run("selecting unknown bank is an error", func(t *testing.T) {
		s := setup(t)

		// s has only 3 banks, so bank 16 is unknown.
		assert.Equal(t, ErrBankUnavailable, s.SelectBank(16))
	})

	t.Run("number of banks must be greater than 1", func(t *testing.T) {
		segment, err := New(0, 0xFF, WithBanks(0))
		assert.Equal(t, ErrBankCountInvalid, err)
		assert.Nil(t, segment)
	})
}

func TestNewWithBanksAndPinnedBank0(t *testing.T) {
	setup := func(t *testing.T) Segment {
		s, err := New(0, 0xFF, WithBanks(3), WithPinnedBank0())
		require.NoError(t, err)
		require.NotNil(t, s)
		return s
	}

	t.Run("contains address works the same as normal segment", func(t *testing.T) {
		s := setup(t)

		assert.True(t, s.ContainsAddress(0x00))
		assert.True(t, s.ContainsAddress(0x80))
		assert.True(t, s.ContainsAddress(0xFF))
		assert.False(t, s.ContainsAddress(0x100))
	})

	t.Run("cannot select bank 0", func(t *testing.T) {
		s := setup(t)
		assert.Equal(t, ErrBankUnavailable, s.SelectBank(0))
	})

	t.Run("default bank is bank 1", func(t *testing.T) {
		s := setup(t)

		ptr := s.BytePtr(0x80)
		require.NotNil(t, ptr)
		assert.Equal(t, uint8(0), *ptr)

		*ptr = 0xff

		require.NoError(t, s.SelectBank(2))

		ptr2 := s.BytePtr(0x80)
		require.NotNil(t, ptr2)
		assertPtrNotEqual(t, ptr, ptr2)
		assert.Equal(t, uint8(0), *ptr2)

		require.NoError(t, s.SelectBank(1))

		ptr3 := s.BytePtr(0x80)
		require.NotNil(t, ptr3)
		assertPtrEqual(t, ptr, ptr3)
		assert.Equal(t, uint8(0xff), *ptr3)
	})

	t.Run("last bank is available", func(t *testing.T) {
		s := setup(t)

		require.NoError(t, s.SelectBank(2))
		require.NotNil(t, s.BytePtr(0x80))
	})

	t.Run("change in bank does not change bank 0", func(t *testing.T) {
		s := setup(t)

		ptr0 := s.BytePtr(0)
		require.NotNil(t, ptr0)

		ptr1 := s.BytePtr(0x80)
		require.NotNil(t, ptr1)

		require.NoError(t, s.SelectBank(2))

		ptr2 := s.BytePtr(0x80)
		assertPtrNotEqual(t, ptr1, ptr2)
		assert.Equal(t, ptr0, s.BytePtr(0))
	})

	t.Run("cannot select a bank outside of bounds", func(t *testing.T) {
		s := setup(t)

		ptr := s.BytePtr(0x80)
		assert.NotNil(t, ptr)

		assert.Equal(t, ErrBankUnavailable, s.SelectBank(3))
		assert.Equal(t, ptr, s.BytePtr(0x80))
	})

	t.Run("cannot pin bank 0 with only one bank", func(t *testing.T) {
		s, err := New(0x00, 0xFF, WithPinnedBank0())
		assert.Nil(t, s)
		assert.Equal(t, ErrCannotPin0WithOneBank, err)
	})
}

func TestNewWithInitialData(t *testing.T) {
	t.Run("initial data is copied into buffer", func(t *testing.T) {
		s, err := New(0x00, 0xFF, WithInitialData([]uint8("\x74\x75\x76\x77")))
		require.NoError(t, err)

		ptr := s.BytePtr(0x02)
		require.NotNil(t, ptr)
		require.Equal(t, uint8(0x76), *ptr)
	})

	t.Run("initial data cannot be bigger than segment", func(t *testing.T) {
		s, err := New(0x00, 0x01, WithInitialData([]uint8("\x74\x75\x76\x77")))
		assert.Nil(t, s)
		assert.Equal(t, ErrBufferIncompatible, err)
	})
}

func TestWithMirrorMapping(t *testing.T) {
	t.Run("mirror mapping maps back to the main segment", func(t *testing.T) {
		s, err := New(0x00, 0xFF, WithMirrorMapping(0x1000, 0x10FF))
		require.NoError(t, err)

		assert.True(t, s.ContainsAddress(0))
		assert.True(t, s.ContainsAddress(0x1000))

		assert.False(t, s.ContainsAddress(0x100))

		ptr1 := s.BytePtr(0x00)
		ptr2 := s.BytePtr(0x1000)
		assertPtrEqual(t, ptr1, ptr2)
	})

	t.Run("cannot use mirroring on smaller segment", func(t *testing.T) {
		s, err := New(0x00, 0xF, WithMirrorMapping(0x1000, 0x10FF))
		assert.Nil(t, s)
		assert.Equal(t, ErrInvalidMirrorRange, err)
	})
}
