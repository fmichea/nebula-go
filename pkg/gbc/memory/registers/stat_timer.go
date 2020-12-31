package registers

var (
	_statModesCycles = map[STATMode]int16{
		STATModeHBlank:        201,
		STATModeVBlank:        456, // will be done 10 times, for 4560 total cycles
		STATModeDataTransfer1: 77,
		STATModeDataTransfer2: 169,
	}
)

type STATTimer struct {
	waitCount int16
}

func NewSTATTimer() *STATTimer {
	return &STATTimer{
		waitCount: 0, // FIXME: was 456 in C++ version but no documentation of why?
	}
}

func (t *STATTimer) Forward(cycles uint16) {
	t.waitCount -= int16(cycles)
}

func (t *STATTimer) Expired() bool {
	return t.waitCount < 0
}

func (t *STATTimer) SwitchMode(mode STATMode) {
	t.waitCount += _statModesCycles[mode]
}
