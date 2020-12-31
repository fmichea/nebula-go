package z80

import (
	"nebula-go/pkg/gbc/memory/registers"
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
	z80lib "nebula-go/pkg/gbc/z80/lib"
)

func (c *CPU) manageInterruptRequests() {
	// Interrupts are enabled only if the CPU flag is on, or we are in HALT mode.
	if !c.Regs.IME && !c.Regs.HaltMode {
		return
	}

	mmuRegs := c.MMU.Registers()

	// There is no requested or enabled interrupt currently.
	if mmuRegs.IF.Get() == 0 || mmuRegs.IE.Get() == 0 {
		return
	}

	interrupts := []struct {
		requestFlag *registers.InterruptFlag
		enableFlag  registerslib.Flag
		int         z80lib.Interrupt
	}{
		{mmuRegs.IF.VBlank, mmuRegs.IE.VBlank, z80lib.Rst40h},
		{mmuRegs.IF.STAT, mmuRegs.IE.STAT, z80lib.Rst48h},
		{mmuRegs.IF.Timer, mmuRegs.IE.Timer, z80lib.Rst50h},
		{mmuRegs.IF.Serial, mmuRegs.IE.Serial, z80lib.Rst58h},
		{mmuRegs.IF.Joypad, mmuRegs.IE.Joypad, z80lib.Rst60h},
	}

	for _, interrupt := range interrupts {
		if !interrupt.requestFlag.IsRequested() {
			continue
		}

		c.Regs.HaltMode = false

		if c.Regs.IME && interrupt.enableFlag.GetBool() {
			// Remove the request for this interrupt.
			interrupt.requestFlag.Acknowledge()

			// Disable interrupts and queue interrupt to be executed.
			c.diOpcode()
			c.inPlaceInterrupts[interrupt.int]()

			// Only first matching interrupt is considered.
			break
		}
	}
}

func (c *CPU) manageTimers(cycles uint16) {
	for x := uint16(0); x < cycles; x += 4 {
		c.MMURegs.DIV.MaybeIncrement(0x4)
		c.MMURegs.TIMA.MaybeIncrement()
	}
}
