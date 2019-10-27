package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbyte_Get(t *testing.T) {
	reg := NewDByte(0xABCD)
	assert.Equal(t, uint16(0xABCD), reg.Get())
}

func TestDbyte_Set(t *testing.T) {
	reg := NewDByte(0xABCD)
	reg.Set(0xCDEF)

	assert.Equal(t, uint16(0xCDEF), reg.Get())
}
