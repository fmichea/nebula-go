package registers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLCDCReg(t *testing.T) {
	var value uint8

	reg := NewLCDCReg(&value)
	require.NotNil(t, reg)
}
