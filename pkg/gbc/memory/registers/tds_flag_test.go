package registers

import (
	"testing"
)

func TestTDSFlag(t *testing.T) {
	//setup := func(value uint8) (registerslib.Byte, *TDSFlag) {
	//	reg := registerslib.NewByte(&value, value)
	//	return reg, NewTDSFlag(reg, 0)
	//}

	t.Run("tile address: unsigned based on 0x8000 when enabled", func(t *testing.T) {
		//_, flag := setup(0xFF)
		//
		//assert.Equal(t, uint16(0x8000), flag.BaseAddress(0))
		//assert.Equal(t, uint16(0x8020), flag.BaseAddress(2))
		//assert.Equal(t, uint16(0x8030), flag.BaseAddress(3))
		//
		//assert.Equal(t, uint16(0x87F0), flag.BaseAddress(0x7F))
		//assert.Equal(t, uint16(0x8800), flag.BaseAddress(0x80))
		//
		//assert.Equal(t, uint16(0x8FB0), flag.BaseAddress(0xFB))
		//assert.Equal(t, uint16(0x8FF0), flag.BaseAddress(0xFF))
	})

	t.Run("tile address: signed based on 0x8800 when disabled", func(t *testing.T) {
		//_, flag := setup(0xFE)
		//
		//assert.Equal(t, uint16(0x9000), flag.BaseAddress(0))
		//assert.Equal(t, uint16(0x9020), flag.BaseAddress(2))
		//assert.Equal(t, uint16(0x9030), flag.BaseAddress(3))
		//
		//assert.Equal(t, uint16(0x97F0), flag.BaseAddress(0x7F))
		//assert.Equal(t, uint16(0x8800), flag.BaseAddress(0x80))
		//
		//assert.Equal(t, uint16(0x8FB0), flag.BaseAddress(0xFB))
		//assert.Equal(t, uint16(0x8FF0), flag.BaseAddress(0xFF))
	})
}
