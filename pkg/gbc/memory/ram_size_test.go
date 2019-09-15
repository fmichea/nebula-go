package memory

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validRAMSizeValues = []RAMSize{
		RAMSizeNone,
		RAMSize2KB,
		RAMSize8KB,
		RAMSize32KB,
	}

	invalidRAMSize = RAMSize(0xff)
)

func TestRAMSize_IsValid(t *testing.T) {

	for _, value := range validRAMSizeValues {
		t.Run(fmt.Sprintf("%s (%x) is valid", value.String(), value), func(t *testing.T) {
			assert.True(t, value.IsValid())
		})
	}

	t.Run("invalid case", func(t *testing.T) {
		assert.False(t, invalidRAMSize.IsValid())
	})
}

func TestRAMSize_String(t *testing.T) {
	for _, value := range validRAMSizeValues {
		t.Run(fmt.Sprintf("%s (%x) is non-empty string", value.String(), value), func(t *testing.T) {
			assert.NotEmpty(t, value.String())
		})
	}

	t.Run("invalid case", func(t *testing.T) {
		assert.Empty(t, invalidRAMSize.String())
	})
}
