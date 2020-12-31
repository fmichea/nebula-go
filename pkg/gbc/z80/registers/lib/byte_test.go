package registerslib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewByte(t *testing.T) {
	reg := NewByte(0x45)
	require.NotNil(t, reg)
	require.Equal(t, uint8(0x45), reg.Get())
}

func TestNewByteWithMask(t *testing.T) {
	reg := NewByteWithMask(0x45, 0xF0)
	require.NotNil(t, reg)
	require.Equal(t, uint8(0x40), reg.Get())
}
