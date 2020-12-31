package registerslib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestByteWithErrors(t *testing.T) {
	var value uint8

	reg1 := NewByte(&value, 0x55)
	reg2 := WrapWithError(reg1)

	val, err := reg2.Get()
	require.NoError(t, err)
	assert.Equal(t, uint8(0x55), val)

	require.NoError(t, reg2.Set(0x88))
	assert.Equal(t, uint8(0x88), value)
}
