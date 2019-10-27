package bitwise

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLowBit8(t *testing.T) {
	cases := []struct {
		initialValue uint8
		bitValue     uint8
	}{
		{0x00, 0},
		{0xFE, 0},
		{0x01, 1},
		{0xFF, 1},

		//{0x00, 7, 0},
		//{0x7F, 7, 0},
		//{0x80, 7, 1},
		//{0xFF, 7, 1},
	}

	for _, c := range cases {
		name := fmt.Sprintf("low bit of %#v is %d", c.initialValue, c.bitValue)

		t.Run(name, func(t *testing.T) {
			assert.Equal(t, c.bitValue, LowBit8(c.initialValue))
		})
	}
}

func TestHighBit8(t *testing.T) {
	cases := []struct {
		initialValue uint8
		bitValue     uint8
	}{
		{0x00, 0},
		{0x7F, 0},
		{0x80, 1},
		{0xFF, 1},
	}

	for _, c := range cases {
		name := fmt.Sprintf("high bit of %#v is %d", c.initialValue, c.bitValue)

		t.Run(name, func(t *testing.T) {
			assert.Equal(t, c.bitValue, HighBit8(c.initialValue))
		})
	}
}
