package z80lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRegisters(t *testing.T) {
	assert.NotNil(t, NewRegisters())
}
