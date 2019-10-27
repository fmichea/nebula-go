package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbyteSplitByte_Get(t *testing.T) {
	reg1 := NewByte(0xAB)
	reg2 := NewByte(0xCD)

	reg := NewSplitDByte(reg1, reg2)
	assert.Equal(t, uint16(0xABCD), reg.Get())
}

func TestDbyteSplitByte_Set(t *testing.T) {
	reg1 := NewByte(0xAB)
	reg2 := NewByte(0xCD)

	reg := NewSplitDByte(reg1, reg2)
	reg.Set(0xBADA)

	assert.Equal(t, uint16(0xBADA), reg.Get())
	assert.Equal(t, uint8(0xBA), reg1.Get())
	assert.Equal(t, uint8(0xDA), reg2.Get())
}
