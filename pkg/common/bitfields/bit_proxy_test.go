package bitfields

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitproxy_Get(t *testing.T) {
	cases := []struct {
		regValue uint8
		bitIdx   uint8
		bitMask  uint8
		bitValue uint8
	}{
		{0x00, 0, 0x1, 0x00},
		{0xFF, 0, 0x1, 0x01},
		{0x01, 0, 0x9, 0x01},
		{0xFF, 0, 0x9, 0x09},
		{0xF0, 4, 0x9, 0x09},
	}

	for _, c := range cases {
		name := fmt.Sprintf(
			"reg with value %#v, but %d (mask = %#v) has value %d",
			c.regValue,
			c.bitIdx,
			c.bitMask,
			c.bitValue,
		)

		t.Run(name, func(t *testing.T) {
			reg := NewByte(c.regValue)
			bitProxy := NewBitProxy(reg, c.bitIdx, c.bitMask)

			assert.Equal(t, c.bitValue, bitProxy.Get())
		})
	}
}

func TestBitproxy_Set(t *testing.T) {
	cases := []struct {
		initialValue uint8
		resultValue  uint8
		bitIdx       uint8
		bitMask      uint8
		bitValue     uint8
	}{
		{0x00, 0x00, 0, 0x1, 0x00},
		{0x00, 0x01, 0, 0x1, 0x01},
		{0x00, 0x10, 4, 0x1, 0x01},
		{0xFF, 0xFF, 0, 0x1, 0x01},
		{0x00, 0x12, 1, 0x9, 0xFF},
		{0xFF, 0xF6, 0, 0x9, 0x00},
		{0xFF, 0x6F, 4, 0x9, 0x00},
	}

	for _, c := range cases {
		name := fmt.Sprintf(
			"reg (initial = %#v, result = %#v), but %d (mask = %#v) has value %d",
			c.initialValue,
			c.resultValue,
			c.bitIdx,
			c.bitMask,
			c.bitValue,
		)

		t.Run(name, func(t *testing.T) {
			reg := NewByte(c.initialValue)

			bitProxy := NewBitProxy(reg, c.bitIdx, c.bitMask)
			bitProxy.Set(c.bitValue)

			assert.Equal(t, c.resultValue, reg.Get())
		})
	}
}
