package z80lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlags_String(t *testing.T) {
	cases := []struct {
		initialFlags uint8
		value        string
	}{
		{FlagsCleared, ""},
		{FlagsFullSet, "ZF | NE | HC | CY"},
		{ZF, "ZF"},
		{ZF | NE, "ZF | NE"},
		{HC | CY, "HC | CY"},
	}

	for _, c := range cases {
		assert.Equal(t, c.value, NewFlags(c.initialFlags).String(), "flags %#v should be %s", c.initialFlags, c.value)
	}
}
