package registers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSTATReg(t *testing.T) {
	var value uint8

	reg := NewSTATReg(&value)
	require.NotNil(t, reg)
}
