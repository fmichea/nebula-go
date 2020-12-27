package clock

import (
	bbjclock "github.com/benbjohnson/clock"
)

type Clock bbjclock.Clock

var New = bbjclock.New
