package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlag_GetBool(t *testing.T) {
	reg := NewByte(0xF0)

	assert.False(t, NewFlag(reg, 0).GetBool())
	assert.True(t, NewFlag(reg, 4).GetBool())
}

func TestFlag_SetBool(t *testing.T) {
	reg := NewByte(0x00)
	f := NewFlag(reg, 0)

	f.SetBool(true)
	assert.Equal(t, uint8(1), reg.Get())

	f.SetBool(false)
	assert.Equal(t, uint8(0), reg.Get())
}
