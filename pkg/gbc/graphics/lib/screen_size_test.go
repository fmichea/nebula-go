package graphicslib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValues(t *testing.T) {
	assert.NotZero(t, Width)
	assert.NotZero(t, Height)
}
