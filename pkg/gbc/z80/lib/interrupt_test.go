package z80lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterrupt_Addr(t *testing.T) {
	assert.Equal(t, uint16(0x18), Rst18h.Addr())
}
