package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDIVReg(t *testing.T) {
	var tac8 uint8

	tac := NewTACReg(&tac8)

	t.Run("DIV reg starts initialized to 0", func(t *testing.T) {
		div := NewDIVReg(nil, tac)

		assert.Equal(t, uint8(0), div.Get())
	})

	t.Run("it takes 256 cycles for DIV increment to happen", func(t *testing.T) {
		div := NewDIVReg(nil, tac)

		assert.Equal(t, uint8(0), div.Get())
		div.MaybeIncrement(0x10)
		assert.Equal(t, uint8(0), div.Get())

		div.MaybeIncrement(0xEF)
		assert.Equal(t, uint8(0), div.Get())

		div.MaybeIncrement(0x01)
		assert.Equal(t, uint8(1), div.Get())

		div.MaybeIncrement(0xFF)
		assert.Equal(t, uint8(1), div.Get())

		div.MaybeIncrement(0x01)
		assert.Equal(t, uint8(2), div.Get())
	})
}
