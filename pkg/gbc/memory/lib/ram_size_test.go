package lib

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
		t.Run(fmt.Sprintf("%s (%#v) is valid", value, value), func(t *testing.T) {
			assert.True(t, value.IsValid())
		})
	}

	t.Run("invalid case", func(t *testing.T) {
		assert.False(t, invalidRAMSize.IsValid())
	})
}

func TestRAMSize_String(t *testing.T) {
	for _, value := range validRAMSizeValues {
		t.Run(fmt.Sprintf("%s (%#v) is non-empty string", value, value), func(t *testing.T) {
			assert.NotEmpty(t, value.String())
		})
	}

	t.Run("invalid case", func(t *testing.T) {
		assert.Empty(t, invalidRAMSize.String())
	})
}

func TestRAMSize_BankCount(t *testing.T) {
	for _, value := range []RAMSize{RAMSizeNone, RAMSize2KB, RAMSize8KB} {
		t.Run(fmt.Sprintf("%s (%#v) is not banked", value, value), func(t *testing.T) {
			assert.Equal(t, uint(1), value.BankCount())
		})
	}

	t.Run(fmt.Sprintf("%s (%#v) is banked in 4 banks of 8KByte", RAMSize32KB, RAMSize32KB), func(t *testing.T) {
		assert.Equal(t, uint(4), RAMSize32KB.BankCount())
	})
}
