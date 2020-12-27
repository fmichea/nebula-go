package clock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClock(t *testing.T) {
	assert.NotEmpty(t, New().Now().String())
}
