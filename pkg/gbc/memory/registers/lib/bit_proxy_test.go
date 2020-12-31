package registerslib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitProxy(t *testing.T) {
	var value uint8

	reg := NewByte(&value, 0x55)
	bp := NewBitProxy(reg, 4, 0xF)

	bp.Set(0)
	assert.Equal(t, uint8(0x05), reg.Get())

	reg.Set(0x87)
	assert.Equal(t, uint8(0x8), bp.Get())
}
