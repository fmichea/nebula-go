package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

func TestObjSizeFlag(t *testing.T) {
	setup := func(value uint8) *ObjSizeFlag {
		reg := registerslib.NewByte(&value, value)
		return NewObjSizeReg(reg, 0)
	}

	assert.Equal(t, int16(8), setup(0xFE).SpriteHeight())
	assert.Equal(t, int16(16), setup(0xFF).SpriteHeight())
}
