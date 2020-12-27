package bitwise

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertHighLow8To16(t *testing.T) {
	assert.Equal(t, uint16(0xABCD), ConvertHighLow8To16(0xAB, 0xCD))
}
