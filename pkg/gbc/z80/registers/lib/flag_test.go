package registerslib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFlag(t *testing.T) {
	reg := NewByte(0xF0)

	f1 := NewFlag(reg, 0)
	require.NotNil(t, f1)
	require.False(t, f1.GetBool())

	f2 := NewFlag(reg, 4)
	require.NotNil(t, f2)
	require.True(t, f2.GetBool())
}
