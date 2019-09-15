package memory

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

	invalidROMSize = ROMSize(0xff)
)

func TestROMSize_IsValid(t *testing.T) {
	for _, value := range validROMSizeValues {
		t.Run(fmt.Sprintf("%s (%x) is valid", value.String(), value), func(t *testing.T) {
			assert.True(t, value.IsValid())
		})
	}

	t.Run("invalid value case", func(t *testing.T) {
		assert.False(t, invalidROMSize.IsValid())
	})
}

func TestROMSize_String(t *testing.T) {
	for _, value := range validROMSizeValues {
		t.Run(fmt.Sprintf("%s (%x) is non-empty string", value.String(), value), func(t *testing.T) {
			assert.NotEmpty(t, value.String())
		})
	}

	t.Run("invalid value case", func(t *testing.T) {
		assert.Empty(t, invalidROMSize.String())
	})
}
