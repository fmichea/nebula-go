package segments

import (
	"nebula-go/pkg/gbc/memory/lib"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	t.Run("bank count and current bank works", func(t *testing.T) {
		s := setup(t)

		assert.Equal(t, uint(0), s.Bank())
		assert.Equal(t, uint(1), s.BankCount())
	})

	t.Run("address ranges", func(t *testing.T) {
		s := setup(t)
		assert.Len(t, s.AddressRanges(), 1)
	})

	t.Run("changes can be read back", func(t *testing.T) {
		s := setup(t)

		value, err := s.ReadByte(0x01)
		require.NoError(t, err)
		assert.Equal(t, uint8(0), value)

		err = s.WriteByte(0x01, 0xFF)
		require.NoError(t, err)

		value2, err := s.ReadByte(0x01)
		require.NoError(t, err)
		assert.Equal(t, uint8(0xff), value2)
	})

	t.Run("invalid write", func(t *testing.T) {
		s := setup(t)

		err := s.WriteByte(0x1000, 0xFF)
		require.Equal(t, lib.ErrInvalidSegmentAddr, err)
	})

	t.Run("can only select bank 0", func(t *testing.T) {
		s := setup(t)

		assert.Equal(t, ErrBankUnavailable, s.SelectBank(1))
		assert.Equal(t, ErrBankUnavailable, s.SelectBank(13))
		assert.NoError(t, s.SelectBank(0))
	})

	t.Run("pointer outside range is nil", func(t *testing.T) {
		s := setup(t)

		_, err := s.ReadByte(0x1000)
		assert.Equal(t, lib.ErrInvalidSegmentAddr, err)
	})

	t.Run("read byte slice", func(t *testing.T) {
		s := setup(t)

		t.Run("simple valid case", func(t *testing.T) {
			values, err := s.ReadByteSlice(0x00, 0x80)
			require.NoError(t, err)
			assert.Len(t, values, 0x80)
		})

		t.Run("address out of range", func(t *testing.T) {
			_, err := s.ReadByteSlice(0x4000, 0x10)
			assert.Equal(t, lib.ErrInvalidSegmentAddr, err)
		})

		t.Run("not enough data in segment is invalid read", func(t *testing.T) {
			_, err := s.ReadByteSlice(0x00, 0x1000)
			assert.Equal(t, ErrSegmentTooSmall, err)
		})
	})

	t.Run("write byte slice", func(t *testing.T) {
		t.Run("simple valid write", func(t *testing.T) {
			s := setup(t)

			err := s.WriteByteSlice(0x00, []uint8{0xAA, 0xBB})
			require.NoError(t, err)

			values, err := s.ReadByteSlice(0x00, 2)
			require.NoError(t, err)
			assert.Equal(t, []uint8{0xAA, 0xBB}, values)
		})

		t.Run("invalid write, buffer too big", func(t *testing.T) {
			s := setup(t)

			err := s.WriteByteSlice(0xFF, []uint8{0xAA, 0xBB})
			require.Equal(t, ErrSegmentTooSmall, err)
		})

		t.Run("valid write accross middle boundary", func(t *testing.T) {
			s := setup(t)

			err := s.WriteByteSlice(0x7F, []uint8{0xAA, 0xBB})
			require.NoError(t, err)

			values, err := s.ReadByteSlice(0x7F, 2)
			require.NoError(t, err)
			assert.Equal(t, []uint8{0xAA, 0xBB}, values)
		})

		t.Run("invalid write, address not known", func(t *testing.T) {
			s := setup(t)

			err := s.WriteByteSlice(0x2000, []uint8{})
			require.Equal(t, lib.ErrInvalidSegmentAddr, err)
		})
	})

	t.Run("byte hook works with no banks", func(t *testing.T) {
		s := setup(t)

		ptr, err := s.ByteHook(0x10)
		require.NoError(t, err)
		assert.NotNil(t, ptr)
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

	t.Run("bank count and current bank works", func(t *testing.T) {
		s := setup(t)

		assert.Equal(t, uint(0), s.Bank())
		assert.Equal(t, uint(3), s.BankCount())
	})

	t.Run("address ranges", func(t *testing.T) {
		s := setup(t)
		assert.Len(t, s.AddressRanges(), 1)
	})

	t.Run("value is the same for a given address without bank change", func(t *testing.T) {
		s := setup(t)

		value, err := s.ReadByte(0x01)
		require.NoError(t, err)
		assert.Equal(t, uint8(0), value)

		err = s.WriteByte(0x01, 0xFF)
		require.NoError(t, err)

		value2, err := s.ReadByte(0x01)
		require.NoError(t, err)
		assert.Equal(t, uint8(0xff), value2)
	})

	t.Run("invalid write", func(t *testing.T) {
		s := setup(t)

		err := s.WriteByte(0x1000, 0xFF)
		require.Equal(t, lib.ErrInvalidSegmentAddr, err)
	})

	t.Run("bank 0 is the default", func(t *testing.T) {
		s := setup(t)

		err := s.WriteByte(0x01, 0xFF)
		require.NoError(t, err)

		assert.NoError(t, s.SelectBank(0))

		value, err := s.ReadByte(0x01)
		require.NoError(t, err)
		assert.Equal(t, uint8(0xff), value)
	})

	t.Run("bank changes change the actual memory values", func(t *testing.T) {
		s := setup(t)

		// by default we should be on bank 0.
		err := s.WriteByte(0x01, 0xFF)
		require.NoError(t, err)

		// moving to bank 1 and checking we did change memory view.
		assert.NoError(t, s.SelectBank(1))

		value, err := s.ReadByte(0x01)
		require.NoError(t, err)
		assert.Equal(t, uint8(0x00), value)

		// switching back to bank 0.
		assert.NoError(t, s.SelectBank(0))

		value2, err := s.ReadByte(0x01)
		require.NoError(t, err)
		assert.Equal(t, uint8(0xFF), value2)
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

	t.Run("read byte slice", func(t *testing.T) {
		s := setup(t)

		t.Run("simple valid case", func(t *testing.T) {
			values, err := s.ReadByteSlice(0x00, 0x80)
			require.NoError(t, err)
			assert.Len(t, values, 0x80)
		})

		t.Run("address out of range", func(t *testing.T) {
			_, err := s.ReadByteSlice(0x4000, 0x10)
			assert.Equal(t, lib.ErrInvalidSegmentAddr, err)
		})

		t.Run("not enough data in segment is invalid read", func(t *testing.T) {
			_, err := s.ReadByteSlice(0x00, 0x1000)
			assert.Equal(t, ErrSegmentTooSmall, err)
		})
	})

	t.Run("write byte slice", func(t *testing.T) {
		t.Run("simple valid write", func(t *testing.T) {
			s := setup(t)

			err := s.WriteByteSlice(0x00, []uint8{0xAA, 0xBB})
			require.NoError(t, err)

			values, err := s.ReadByteSlice(0x00, 2)
			require.NoError(t, err)
			assert.Equal(t, []uint8{0xAA, 0xBB}, values)
		})

		t.Run("invalid write, buffer too big", func(t *testing.T) {
			s := setup(t)

			err := s.WriteByteSlice(0xFF, []uint8{0xAA, 0xBB})
			require.Equal(t, ErrSegmentTooSmall, err)
		})

		t.Run("valid write accross middle boundary", func(t *testing.T) {
			s := setup(t)

			err := s.WriteByteSlice(0x7F, []uint8{0xAA, 0xBB})
			require.NoError(t, err)

			values, err := s.ReadByteSlice(0x7F, 2)
			require.NoError(t, err)
			assert.Equal(t, []uint8{0xAA, 0xBB}, values)
		})

		t.Run("invalid write, address not known", func(t *testing.T) {
			s := setup(t)

			err := s.WriteByteSlice(0x2000, []uint8{})
			require.Equal(t, lib.ErrInvalidSegmentAddr, err)
		})
	})

	t.Run("byte hook is invalid anywhere", func(t *testing.T) {
		s := setup(t)

		_, err := s.ByteHook(0x10)
		assert.Equal(t, ErrInvalidHookInBank, err)
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

	t.Run("bank count and current bank works", func(t *testing.T) {
		s := setup(t)

		assert.Equal(t, uint(1), s.Bank())
		assert.Equal(t, uint(3), s.BankCount())
	})

	t.Run("address ranges", func(t *testing.T) {
		s := setup(t)
		assert.Len(t, s.AddressRanges(), 1)
	})

	t.Run("select bank 0 is bank 1", func(t *testing.T) {
		s := setup(t)

		// Select bank 1 and write 0xFF at address.
		assert.NoError(t, s.SelectBank(1))

		err := s.WriteByte(0x80, 0xFF)
		require.NoError(t, err)

		// Select bank 2 and check the change worked (value == 0x00)
		require.NoError(t, s.SelectBank(2))

		value, err := s.ReadByte(0x80)
		require.NoError(t, err)
		assert.Equal(t, uint8(0x00), value)

		// Select bank 0, which should actually select bank 1.
		require.NoError(t, s.SelectBank(0))

		value2, err := s.ReadByte(0x80)
		require.NoError(t, err)
		assert.Equal(t, uint8(0xFF), value2)
	})

	t.Run("invalid write", func(t *testing.T) {
		s := setup(t)

		err := s.WriteByte(0x1000, 0xFF)
		require.Equal(t, lib.ErrInvalidSegmentAddr, err)
	})

	t.Run("default bank is bank 1", func(t *testing.T) {
		s := setup(t)

		// Default bank should be bank 1.
		err := s.WriteByte(0x80, 0xFF)
		require.NoError(t, err)

		// Select bank 2 and check the change worked (value == 0x00)
		require.NoError(t, s.SelectBank(2))

		value, err := s.ReadByte(0x80)
		require.NoError(t, err)
		assert.Equal(t, uint8(0x00), value)

		// Select bank 1 again.
		require.NoError(t, s.SelectBank(1))

		value2, err := s.ReadByte(0x80)
		require.NoError(t, err)
		assert.Equal(t, uint8(0xFF), value2)
	})

	t.Run("last bank is available", func(t *testing.T) {
		s := setup(t)

		// There are 3 banks, so bank 2 is the last one.
		require.NoError(t, s.SelectBank(2))

		value, err := s.ReadByte(0x80)
		require.NoError(t, err)
		assert.Equal(t, uint8(0x00), value)
	})

	t.Run("change in bank does not change bank 0", func(t *testing.T) {
		s := setup(t)

		// Write in bank 0 and 1.
		err := s.WriteByte(0x10, 0x88)
		require.NoError(t, err)

		err = s.WriteByte(0x80, 0xBB)
		require.NoError(t, err)

		// Switch to bank 2.
		require.NoError(t, s.SelectBank(2))

		value, err := s.ReadByte(0x10)
		require.NoError(t, err)
		assert.Equal(t, uint8(0x88), value)

		value2, err := s.ReadByte(0x80)
		require.NoError(t, err)
		assert.Equal(t, uint8(0x00), value2)
	})

	t.Run("cannot select a bank outside of bounds", func(t *testing.T) {
		s := setup(t)

		assert.Equal(t, ErrBankUnavailable, s.SelectBank(3))
	})

	t.Run("cannot pin bank 0 with only one bank", func(t *testing.T) {
		s, err := New(0x00, 0xFF, WithPinnedBank0())
		assert.Nil(t, s)
		assert.Equal(t, ErrCannotPin0WithOneBank, err)
	})

	t.Run("read byte slice", func(t *testing.T) {
		s := setup(t)

		t.Run("simple valid case, no cross over banks boundaries", func(t *testing.T) {
			values, err := s.ReadByteSlice(0x00, 0x20)
			require.NoError(t, err)
			assert.Len(t, values, 0x20)
		})

		t.Run("valid case with cross over banks boundaries", func(t *testing.T) {
			values, err := s.ReadByteSlice(0x70, 0x20)
			require.NoError(t, err)
			assert.Len(t, values, 0x20)
		})

		t.Run("address out of range", func(t *testing.T) {
			_, err := s.ReadByteSlice(0x4000, 0x10)
			assert.Equal(t, lib.ErrInvalidSegmentAddr, err)
		})

		t.Run("not enough data in segment is invalid read", func(t *testing.T) {
			_, err := s.ReadByteSlice(0x00, 0x1000)
			assert.Equal(t, ErrSegmentTooSmall, err)
		})
	})

	t.Run("write byte slice", func(t *testing.T) {
		t.Run("simple valid write", func(t *testing.T) {
			s := setup(t)

			err := s.WriteByteSlice(0x00, []uint8{0xAA, 0xBB})
			require.NoError(t, err)

			values, err := s.ReadByteSlice(0x00, 2)
			require.NoError(t, err)
			assert.Equal(t, []uint8{0xAA, 0xBB}, values)
		})

		t.Run("invalid write, buffer too big", func(t *testing.T) {
			s := setup(t)

			err := s.WriteByteSlice(0xFF, []uint8{0xAA, 0xBB})
			require.Equal(t, ErrSegmentTooSmall, err)
		})

		t.Run("valid write accross middle boundary", func(t *testing.T) {
			s := setup(t)

			err := s.WriteByteSlice(0x7F, []uint8{0xAA, 0xBB})
			require.NoError(t, err)

			values, err := s.ReadByteSlice(0x7F, 2)
			require.NoError(t, err)
			assert.Equal(t, []uint8{0xAA, 0xBB}, values)
		})

		t.Run("invalid write, address not known", func(t *testing.T) {
			s := setup(t)

			err := s.WriteByteSlice(0x2000, []uint8{})
			require.Equal(t, lib.ErrInvalidSegmentAddr, err)
		})
	})

	t.Run("byte hook only valid in bank 0 (first half)", func(t *testing.T) {
		s := setup(t)

		ptr, err := s.ByteHook(0x10)
		require.NoError(t, err)
		assert.NotNil(t, ptr)

		_, err = s.ByteHook(0x80)
		assert.Equal(t, ErrInvalidHookInBank, err)
	})
}

func TestNewWithInitialData(t *testing.T) {
	t.Run("initial data is copied into buffer", func(t *testing.T) {
		s, err := New(0x00, 0xFF, WithInitialData([]uint8("\x74\x75\x76\x77")))
		require.NoError(t, err)

		value, err := s.ReadByte(0x02)
		require.NoError(t, err)
		require.Equal(t, uint8(0x76), value)
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

		assert.Len(t, s.AddressRanges(), 2)

		err = s.WriteByte(0x00, 0xFF)
		require.NoError(t, err)

		value, err := s.ReadByte(0x00)
		require.NoError(t, err)
		require.Equal(t, uint8(0xFF), value)

		value2, err := s.ReadByte(0x1000)
		require.NoError(t, err)
		require.Equal(t, uint8(0xFF), value2)
	})

	t.Run("cannot use mirroring on smaller segment", func(t *testing.T) {
		s, err := New(0x00, 0xF, WithMirrorMapping(0x1000, 0x10FF))
		assert.Nil(t, s)
		assert.Equal(t, ErrInvalidMirrorRange, err)
	})
}
