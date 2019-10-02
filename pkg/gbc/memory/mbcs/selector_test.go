package mbcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelector_Name(t *testing.T) {
	assert.Equal(t, "ROM ONLY", NewSelector(0x00).Name())
	assert.Equal(t, "MBC4+RAM", NewSelector(0x16).Name())
	assert.Equal(t, "POCKET CAMERA", NewSelector(0xFC).Name())
	assert.Equal(t, "BANDAI TAMA5", NewSelector(0xFD).Name())
	assert.Equal(t, "HuC3", NewSelector(0xFE).Name())
	assert.Equal(t, "HuC1+RAM+BATTERY", NewSelector(0xFF).Name())
}

func TestSelector_IsValid(t *testing.T) {
	assert.True(t, NewSelector(0x00).IsValid())
	assert.True(t, NewSelector(0x16).IsValid())
	assert.True(t, NewSelector(0xFC).IsValid())
	assert.True(t, NewSelector(0xFD).IsValid())
	assert.True(t, NewSelector(0xFF).IsValid())

	assert.False(t, NewSelector(0x80).IsValid())
}

func TestSelector_GetMBC(t *testing.T) {
	assert.NotNil(t, NewSelector(0x00).GetMBC(nil, nil))
	assert.NotNil(t, NewSelector(0x01).GetMBC(nil, nil))

	// FIXME: implement the other MBCs
	assert.Nil(t, NewSelector(0x05).GetMBC(nil, nil))
	assert.Nil(t, NewSelector(0x0F).GetMBC(nil, nil))
	assert.Nil(t, NewSelector(0x1A).GetMBC(nil, nil))

	assert.Nil(t, NewSelector(0xFF).GetMBC(nil, nil))
}
