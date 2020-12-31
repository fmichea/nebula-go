package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

func TestSTATModeBitProxy(t *testing.T) {
	var value uint8

	reg := registerslib.NewByte(&value, value)
	timer := NewSTATTimer()

	bitproxy := NewSTATModeBitProxy(reg, timer)

	assert.Equal(t, STATModeHBlank, bitproxy.GetMode())
	assert.Equal(t, int16(0), timer.waitCount)

	bitproxy.SetMode(STATModeVBlank)
	assert.Equal(t, STATModeVBlank, bitproxy.GetMode())
	assert.NotEqual(t, int16(0), timer.waitCount)
}
