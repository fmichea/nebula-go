package bitwise

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBits8(t *testing.T) {
	cases := []struct {
		value  uint8
		offset uint8
		mask   uint8
		result uint8
	}{
		{0x00, 0, 0x00, 0x00},
		{0x00, 0, 0xFF, 0x00},
		{0xFF, 0, 0x00, 0x00},
		{0x0F, 0, 0x0F, 0x0F},
		{0xF0, 4, 0x0F, 0x0F},
		{0xF0, 3, 0x05, 0x04},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("TestGetBits8: %v", c), func(t *testing.T) {
			assert.Equal(t, c.result, GetBits8(c.value, c.offset, c.mask))
		})
	}
}

func TestGetBit8(t *testing.T) {
	cases := []struct {
		value  uint8
		offset uint8
		result uint8
	}{
		{0x00, 0, 0x00},
		{0xFF, 0, 0x01},
		{0x0F, 0, 0x01},
		{0xF0, 4, 0x01},
		{0xF0, 3, 0x00},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("TestGetBit8: %v", c), func(t *testing.T) {
			assert.Equal(t, c.result, GetBit8(c.value, c.offset))
		})
	}
}

func TestGetBits16(t *testing.T) {
	cases := []struct {
		value  uint16
		offset uint16
		mask   uint16
		result uint16
	}{
		{0x0000, 0, 0x0000, 0x0000},
		{0x0000, 0, 0x00FF, 0x0000},
		{0x00FF, 0, 0x0000, 0x0000},
		{0x000F, 0, 0x000F, 0x000F},
		{0x00F0, 4, 0x000F, 0x000F},
		{0x00F0, 3, 0x0005, 0x0004},
		{0x00F0, 8, 0x0005, 0x0000},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("TestGetBits16: %v", c), func(t *testing.T) {
			assert.Equal(t, c.result, GetBits16(c.value, c.offset, c.mask))
		})
	}
}

func TestGetBit16(t *testing.T) {
	cases := []struct {
		value  uint16
		offset uint16
		result uint16
	}{
		{0x00, 0, 0x00},
		{0xFF, 0, 0x01},
		{0x0F, 0, 0x01},
		{0xF0, 4, 0x01},
		{0xF0, 3, 0x00},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("TestGetBit16: %v", c), func(t *testing.T) {
			assert.Equal(t, c.result, GetBit16(c.value, c.offset))
		})
	}
}
