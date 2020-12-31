package registerslib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlag(t *testing.T) {
	var value uint8

	reg := NewByte(&value, 0x55)

	flag := NewFlag(reg, 0)
	assert.True(t, flag.GetBool())

	flag.SetBool(false)
	assert.Equal(t, uint8(0x54), reg.Get())

	flag.SetBool(true)
	assert.Equal(t, uint8(0x55), reg.Get())
}
