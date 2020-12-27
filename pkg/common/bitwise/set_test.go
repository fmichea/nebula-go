package bitwise

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetBits8(t *testing.T) {
	cases := []struct {
		current uint8
		offset  uint8
		mask    uint8
		value   uint8
		result  uint8
	}{
		{0x00, 0, 0x00, 0x00, 0x00},
		{0x00, 0, 0x11, 0x0F, 0x01},
		{0x00, 0, 0x11, 0xFF, 0x11},
		{0x00, 4, 0x03, 0xFF, 0x30},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("TestSetBits8: %v", c), func(t *testing.T) {
			assert.Equal(t, c.result, SetBits8(c.current, c.offset, c.mask, c.value))
		})
	}
}
