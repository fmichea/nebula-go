package joypad

import (
	"nebula-go/pkg/common/frontends"
)

func (s *unitTestSuite) TestEscapePressedStopsCPU() {
	s.False(s.mockRegs.Stopped)
	s.callback(frontends.EscapeKey, frontends.DownState)
	s.True(s.mockRegs.Stopped)
}

func (s *unitTestSuite) TestButtonPressedSetsRegsProperly() {
	s.Equal(uint8(0xFF), s.mockRegs.JOYP.Get())
	s.False(s.mockRegs.IF.Joypad.GetBool())

	// Press A, sets the right button and requests interrupt.
	s.callback(frontends.AKey, frontends.DownState)
	s.False(s.mockRegs.JOYP.AButton.GetBool())
	s.True(s.mockRegs.IF.Joypad.GetBool())

	// Press Down, also sets the right button, A still pressed.
	s.callback(frontends.DownKey, frontends.DownState)
	s.False(s.mockRegs.JOYP.AButton.GetBool())
	s.False(s.mockRegs.JOYP.DownButton.GetBool())
	s.True(s.mockRegs.IF.Joypad.GetBool())

	// Unpress Down, unsets that button but A still pressed.
	s.callback(frontends.DownKey, frontends.UpState)
	s.False(s.mockRegs.JOYP.AButton.GetBool())
	s.True(s.mockRegs.JOYP.DownButton.GetBool())
	s.True(s.mockRegs.IF.Joypad.GetBool())
}

func (s *unitTestSuite) TestUnknownKeyIsIgnored() {
	s.callback(frontends.KeyboardKey(55), frontends.DownState)
	s.Equal(uint8(0xFF), s.mockRegs.JOYP.Get())
}
