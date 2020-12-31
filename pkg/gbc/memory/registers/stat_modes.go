package registers

type STATMode int

const (
	STATModeHBlank        STATMode = iota // H-Blank
	STATModeVBlank                        // V-Blank
	STATModeDataTransfer1                 // OAM locked
	STATModeDataTransfer2                 // OAM+VRAM locked
)
