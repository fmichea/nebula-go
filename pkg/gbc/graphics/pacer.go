package graphics

import (
	"fmt"
	"time"

	"nebula-go/pkg/common/clock"
)

const (
	frameDuration = 16600000 * time.Nanosecond
)

type Pacer struct {
	lastTick time.Time
	clock    clock.Clock

	frameRateSum time.Duration
	count        time.Duration
}

func NewPacer(clock clock.Clock) *Pacer {
	return &Pacer{
		lastTick: clock.Now(),
		clock:    clock,
	}
}

func (p *Pacer) Wait() {
	since := p.clock.Since(p.lastTick)
	p.computeFrameRateAverage(since)

	waitTime := frameDuration - since
	if 0 < waitTime {
		p.clock.Sleep(waitTime)
	}
	p.lastTick = p.clock.Now()
}

func (p *Pacer) computeFrameRateAverage(d time.Duration) {
	p.frameRateSum += d
	p.count++

	if p.count == 60 {
		fmt.Println("Frame generation time (per second):", p.frameRateSum/p.count)

		p.frameRateSum = 0
		p.count = 0
	}
}
