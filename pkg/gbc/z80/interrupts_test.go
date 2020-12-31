package z80

func (s *unitTestSuite) TestInterruptsIMEFalseNotHALTIsNoop() {
	// We prepare for joypad interrupt to be requested to ensure that is not
	// executed.
	s.mockRegs.IF.Joypad.SetBool(true)
	s.mockRegs.IE.Joypad.SetBool(true)

	// Disable interrupts and launch manager for this cycle.
	s.cpu.Regs.IME = false
	s.cpu.manageInterruptRequests()

	// Joypad interrupt still requested.
	s.True(s.mockRegs.IF.Joypad.GetBool())
}

func (s *unitTestSuite) TestInterruptsIMEFalseButHALTIsChecked() {
	// We prepare for joypad interrupt to be requested to ensure that is not
	// executed but we get out of halt mode.
	s.mockRegs.IF.Joypad.SetBool(true)
	s.mockRegs.IE.Joypad.SetBool(true)

	// Disable interrupts, but set it in halt mode. Interrupts will be executed.
	s.cpu.Regs.IME = false
	s.cpu.Regs.HaltMode = true

	s.cpu.manageInterruptRequests()

	// Joypad interrupt is still requested. That is because having an interrupt
	// requested just gets us out of halt mode, but does not execute the
	// interrupt yet.
	// FIXME: This is the implementation, I seem to remember this is due to the
	//  fact that the real hardware executes one instruction right after HALT
	//  mode but before the interrupt. But need to be verified.
	s.True(s.mockRegs.IF.Joypad.GetBool())
	s.False(s.cpu.Regs.HaltMode)
}

func (s *unitTestSuite) TestInterruptsNotRequestedIsNoop() {
	s.cpu.Regs.IME = true
	s.cpu.manageInterruptRequests()
}

func (s *unitTestSuite) TestInterruptRequestedAndEnabledIsExecuted() {
	sp := s.cpu.Regs.SP

	s.cpu.Regs.IME = true

	// We prepare for joypad interrupt to be requested to ensure that will be executed.
	s.mockRegs.IF.Joypad.SetBool(true)
	s.mockRegs.IE.Joypad.SetBool(true)

	// Expect exec of interrupt to be started by executing a call function.
	s.mockMMU.EXPECT().WriteDByte(sp.Get()-2, uint16(0x100)).Return(nil)

	s.cpu.manageInterruptRequests()

	// Interrupt is not requested anymore, and was executed.
	s.False(s.mockRegs.IF.Joypad.GetBool())
}
