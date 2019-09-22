package lib

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validROMSizeValues = []ROMSize{
		ROMSize32KB,
		ROMSize64KB,
		ROMSize128KB,
		ROMSize256KB,
		ROMSize512KB,
		ROMSize1MB,
		ROMSize2MB,
		ROMSize4MB,
		ROMSize1p1MB,
		ROMSize1p2MB,
		ROMSize1p5MB,
	}

	invalidROMSize = ROMSize(0xFF)
)

func TestROMSize_IsValid(t *testing.T) {
	for _, value := range validROMSizeValues {
		t.Run(fmt.Sprintf("%s (%#v) is valid", value, value), func(t *testing.T) {
			assert.True(t, value.IsValid())
		})
	}

	t.Run("invalid value case", func(t *testing.T) {
		assert.False(t, invalidROMSize.IsValid())
	})
}

func TestROMSize_String(t *testing.T) {
	for _, value := range validROMSizeValues {
		t.Run(fmt.Sprintf("%s (%#v) is non-empty string", value, value), func(t *testing.T) {
			assert.NotEmpty(t, value.String())
		})
	}

	t.Run("invalid value case", func(t *testing.T) {
		assert.Empty(t, invalidROMSize.String())
	})
}

func TestROMSize_BankCount(t *testing.T) {
	values := map[ROMSize]uint{
		ROMSize32KB:  2,
		ROMSize64KB:  4,
		ROMSize128KB: 8,
		ROMSize256KB: 16,
		ROMSize512KB: 32,
		ROMSize1MB:   64,
		ROMSize2MB:   128,
		ROMSize4MB:   256,
		ROMSize1p1MB: 72,
		ROMSize1p2MB: 80,
		ROMSize1p5MB: 96,
	}

	for size, bankCount := range values {
		t.Run(fmt.Sprintf("%s (%#v) has %d banks", size, size, bankCount), func(t *testing.T) {
			assert.Equal(t, bankCount, size.BankCount())
		})
	}

	t.Run("invalid value has 2 banks anyway", func(t *testing.T) {
		assert.Equal(t, uint(2), invalidROMSize.BankCount())
	})
}
