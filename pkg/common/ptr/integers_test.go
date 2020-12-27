package ptr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUInt8(t *testing.T) {
	val := UInt8(10)
	require.NotNil(t, val)
	require.Equal(t, uint8(10), *val)
}
