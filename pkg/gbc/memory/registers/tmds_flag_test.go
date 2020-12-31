package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

func TestTMDSFlag(t *testing.T) {
	setup := func(value uint8) (registerslib.Byte, *TMDSFlag) {
		reg := registerslib.NewByte(&value, value)
		return reg, NewTMDSFlag(reg, 0)
	}

	t.Run("base address is 0x9C00 when TMDS is set", func(t *testing.T) {
		reg, flag := setup(0x00)

		reg.Set(0xFF)
		assert.Equal(t, uint16(0x9C00), flag.BaseAddress())
	})

	t.Run("base address is 0x9800 when TMDS is unset", func(t *testing.T) {
		reg, flag := setup(0x00)

		reg.Set(0xFE)
		assert.Equal(t, uint16(0x9800), flag.BaseAddress())
	})
}
