package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRegisters(t *testing.T) {
	assert.NotNil(t, New())
}
