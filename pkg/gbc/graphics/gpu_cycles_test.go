package graphics

import (
	"time"

	"github.com/golang/mock/gomock"

	"nebula-go/pkg/common/testhelpers"
	"nebula-go/pkg/gbc/memory/lib"
	"nebula-go/pkg/gbc/memory/registers"
)

func (s *unitTestSuite) TestDoHBlank() {
	s.Run("DMG: basic case with drawing of the line", func() {
		s.tctx.cr.Type = lib.DMG01

		s.tctx.mmuRegs.STAT.Mode.SetMode(registers.STATModeHBlank)

		// Enable on screen displaying for now, but disable everything on it
		// to reduce mmu mocking. The drawing function is tested elsewhere.
		s.tctx.mmuRegs.LCDC.LDE.SetBool(true)
		s.tctx.mmuRegs.LCDC.BGD.SetBool(false)
		s.tctx.mmuRegs.LCDC.WDE.SetBool(false)
		s.tctx.mmuRegs.LCDC.OBJSDE.SetBool(false)
		s.tctx.mockWindow.EXPECT().DrawLine(uint(0x10), gomock.Any()).Return(nil)

		// We to prepare for all the side effects of H-Blank.
		s.tctx.mmuRegs.LY.Set(0x10)
		s.tctx.mmuRegs.LYC.Set(0x55)

		s.tctx.mmuRegs.IF.Set(0x00)

		// This flag will be false once the H-Blank is executed.
		s.tctx.mmuRegs.STAT.CoincidenceFlag.SetBool(true)

		// HBlank lasts between 201 and 207 clock cycles, this one will not expire
		// the timer.
		s.Require().NoError(s.tctx.gpu.DoCycles(200))
		s.Equal(registers.STATModeHBlank, s.tctx.mmuRegs.STAT.Mode.GetMode())
		s.True(s.tctx.mmuRegs.STAT.CoincidenceFlag.GetBool())

		// Now we will trigger the execution of the H-Blank and let it go into the
		// next state. In this case, it will be DataTransfer1 (mode 2).
		s.Require().NoError(s.tctx.gpu.DoCycles(10))

		s.Equal(registers.STATModeDataTransfer1, s.tctx.mmuRegs.STAT.Mode.GetMode())
		s.False(s.tctx.mmuRegs.STAT.CoincidenceFlag.GetBool())
		s.False(s.tctx.mmuRegs.IF.STAT.GetBool())
	})

	s.Run("DMG: drawing disabled makes it skipped entirely", func() {
		s.tctx.cr.Type = lib.DMG01

		s.tctx.mmuRegs.STAT.Mode.SetMode(registers.STATModeHBlank)

		// Disable LCD rendering completely
		s.tctx.mmuRegs.LCDC.LDE.SetBool(false)

		// We to prepare for all the side effects of H-Blank.
		s.tctx.mmuRegs.LY.Set(0x10)
		s.tctx.mmuRegs.LYC.Set(0x55)

		s.tctx.mmuRegs.IF.Set(0x00)

		// This flag will be false once the H-Blank is executed.
		s.tctx.mmuRegs.STAT.CoincidenceFlag.SetBool(true)

		// Now we will trigger the execution of the H-Blank and let it go into the
		// next state. In this case, it will be DataTransfer1 (mode 2).
		s.Require().NoError(s.tctx.gpu.DoCycles(210))

		s.Equal(registers.STATModeDataTransfer1, s.tctx.mmuRegs.STAT.Mode.GetMode())
		s.False(s.tctx.mmuRegs.STAT.CoincidenceFlag.GetBool())
		s.False(s.tctx.mmuRegs.IF.STAT.GetBool())
	})

	s.Run("DMG: rendering error is overall error", func() {
		s.tctx.cr.Type = lib.DMG01

		s.tctx.mmuRegs.STAT.Mode.SetMode(registers.STATModeHBlank)

		// Enable on screen displaying for now, but disable everything on it
		// to reduce mmu mocking. The drawing function is tested elsewhere.
		s.tctx.mmuRegs.LCDC.LDE.SetBool(true)
		s.tctx.mmuRegs.LCDC.BGD.SetBool(false)
		s.tctx.mmuRegs.LCDC.WDE.SetBool(false)
		s.tctx.mmuRegs.LCDC.OBJSDE.SetBool(false)
		s.tctx.mockWindow.EXPECT().DrawLine(uint(0x10), gomock.Any()).Return(testhelpers.ErrTesting1)

		// We to prepare for all the side effects of H-Blank.
		s.tctx.mmuRegs.LY.Set(0x10)
		s.tctx.mmuRegs.LYC.Set(0x55)

		s.tctx.mmuRegs.IF.Set(0x00)

		// This flag will be false once the H-Blank is executed.
		s.tctx.mmuRegs.STAT.CoincidenceFlag.SetBool(true)

		// HBlank lasts between 201 and 207 clock cycles, this one will not expire
		// the timer.
		s.Require().NoError(s.tctx.gpu.DoCycles(200))
		s.Equal(registers.STATModeHBlank, s.tctx.mmuRegs.STAT.Mode.GetMode())
		s.True(s.tctx.mmuRegs.STAT.CoincidenceFlag.GetBool())

		// Now we will trigger the execution of the H-Blank and let it go into the
		// next state. In this case, it will be DataTransfer1 (mode 2).
		err := s.tctx.gpu.DoCycles(210)
		s.Require().Equal(testhelpers.ErrTesting1, err)
	})

	s.Run("DMG: coincidence flag set and interrupt requested when same line", func() {
		s.tctx.cr.Type = lib.DMG01

		s.tctx.mmuRegs.STAT.Mode.SetMode(registers.STATModeHBlank)

		// Disable LCD rendering completely
		s.tctx.mmuRegs.LCDC.LDE.SetBool(false)

		// We to prepare for all the side effects of H-Blank.
		s.tctx.mmuRegs.LY.Set(0x55)
		s.tctx.mmuRegs.LYC.Set(0x55)

		s.tctx.mmuRegs.IF.Set(0x00)

		// This flag will be true once the H-Blank is executed.
		s.tctx.mmuRegs.STAT.CoincidenceFlag.SetBool(false)

		// Now we will trigger the execution of the H-Blank and let it go into the
		// next state. In this case, it will be DataTransfer1 (mode 2).
		s.Require().NoError(s.tctx.gpu.DoCycles(210))

		s.Equal(registers.STATModeDataTransfer1, s.tctx.mmuRegs.STAT.Mode.GetMode())
		s.True(s.tctx.mmuRegs.STAT.CoincidenceFlag.GetBool())
		s.True(s.tctx.mmuRegs.IF.STAT.GetBool())
	})

	s.Run("DMG: last line triggers V-Blank", func() {
		s.tctx.cr.Type = lib.DMG01

		s.tctx.mmuRegs.STAT.Mode.SetMode(registers.STATModeHBlank)

		// Disable LCD rendering completely
		s.tctx.mmuRegs.LCDC.LDE.SetBool(false)

		// We to prepare for all the side effects of H-Blank.
		s.tctx.mmuRegs.LY.Set(143)
		s.tctx.mmuRegs.LYC.Set(0x00)

		s.tctx.mmuRegs.IF.Set(0x00)

		// This flag will be false once the H-Blank is executed.
		s.tctx.mmuRegs.STAT.CoincidenceFlag.SetBool(true)

		// Now we will trigger the execution of the H-Blank and let it go into the
		// next state. In this case, it will be DataTransfer1 (mode 2).
		s.Require().NoError(s.tctx.gpu.DoCycles(210))

		s.Equal(registers.STATModeVBlank, s.tctx.mmuRegs.STAT.Mode.GetMode())
		s.False(s.tctx.mmuRegs.STAT.CoincidenceFlag.GetBool())
		s.False(s.tctx.mmuRegs.IF.STAT.GetBool())
	})

	s.Run("CGB: HDMA transfer inactive doesnt trigger it", func() {
		s.tctx.cr.Type = lib.CGB001

		s.tctx.mmuRegs.STAT.Mode.SetMode(registers.STATModeHBlank)

		// Disable HDMA transfer.
		s.tctx.mockHDMA5.EXPECT().MaybeDoHDMA().Return(nil)

		// Disable LCD rendering completely
		s.tctx.mmuRegs.LCDC.LDE.SetBool(false)

		// We to prepare for all the side effects of H-Blank.
		s.tctx.mmuRegs.LY.Set(0x10)
		s.tctx.mmuRegs.LYC.Set(0x55)

		s.tctx.mmuRegs.IF.Set(0x00)

		// This flag will be false once the H-Blank is executed.
		s.tctx.mmuRegs.STAT.CoincidenceFlag.SetBool(true)

		// Now we will trigger the execution of the H-Blank and let it go into the
		// next state. In this case, it will be DataTransfer1 (mode 2).
		s.Require().NoError(s.tctx.gpu.DoCycles(210))

		s.Equal(registers.STATModeDataTransfer1, s.tctx.mmuRegs.STAT.Mode.GetMode())
		s.False(s.tctx.mmuRegs.STAT.CoincidenceFlag.GetBool())
		s.False(s.tctx.mmuRegs.IF.STAT.GetBool())
	})

	s.Run("CGB: HDMA transfer active gets triggered", func() {
		s.tctx.cr.Type = lib.CGB001

		s.tctx.mmuRegs.STAT.Mode.SetMode(registers.STATModeHBlank)

		// HDMA transfer.
		s.tctx.mockHDMA5.EXPECT().MaybeDoHDMA().Return(nil)

		// Disable LCD rendering completely
		s.tctx.mmuRegs.LCDC.LDE.SetBool(false)

		// We to prepare for all the side effects of H-Blank.
		s.tctx.mmuRegs.LY.Set(0x10)
		s.tctx.mmuRegs.LYC.Set(0x55)

		s.tctx.mmuRegs.IF.Set(0x00)

		// This flag will be false once the H-Blank is executed.
		s.tctx.mmuRegs.STAT.CoincidenceFlag.SetBool(true)

		// Now we will trigger the execution of the H-Blank and let it go into the
		// next state. In this case, it will be DataTransfer1 (mode 2).
		s.Require().NoError(s.tctx.gpu.DoCycles(210))

		s.Equal(registers.STATModeDataTransfer1, s.tctx.mmuRegs.STAT.Mode.GetMode())
		s.False(s.tctx.mmuRegs.STAT.CoincidenceFlag.GetBool())
		s.False(s.tctx.mmuRegs.IF.STAT.GetBool())
	})
}

func (s *unitTestSuite) TestDoVBlank() {
	s.Run("return to V-Blank when line is below max LY", func() {
		// Set the current line to line 153.
		s.tctx.mmuRegs.LY.Set(19)

		// Get into V-Blank mode.
		s.tctx.mmuRegs.STAT.Mode.SetMode(registers.STATModeVBlank)

		// Setup so interrupt is currently not requested.
		s.tctx.mmuRegs.IF.Set(0x00)

		// Trigger the V-Blank using up all the cycles.
		s.Require().NoError(s.tctx.gpu.DoCycles(460))

		// Check LY got incremented but we are still in V-Blank.
		s.Equal(uint8(20), s.tctx.mmuRegs.LY.Get())
		s.Equal(registers.STATModeVBlank, s.tctx.mmuRegs.STAT.Mode.GetMode())
	})

	s.Run("switch to data transfer 1 when we reach line 153", func() {
		// Set the current line to line 153.
		s.tctx.mmuRegs.LY.Set(_maxLY)

		// Get into V-Blank mode.
		s.tctx.mmuRegs.STAT.Mode.SetMode(registers.STATModeVBlank)

		// Setup so interrupt is currently not requested.
		s.tctx.mmuRegs.IF.Set(0x00)

		// We pacer to compute the sleep time, sleep then fetch the new time.
		s.tctx.mockClock.EXPECT().Since(s.tctx.originalTime).Return(1400000 * time.Nanosecond)
		s.tctx.mockClock.EXPECT().Sleep(15200000 * time.Nanosecond)
		s.tctx.mockClock.EXPECT().Now().Return(s.tctx.originalTime.Add(17 * time.Millisecond))

		// The lines drawn will be committed.
		s.tctx.mockWindow.EXPECT().Commit().Return(nil)

		// Trigger the V-Blank using up all the cycles. This will request interrupt, reset LY and get into
		// data transfer mode 1.
		s.Require().NoError(s.tctx.gpu.DoCycles(460))

		s.Equal(registers.STATModeDataTransfer1, s.tctx.mmuRegs.STAT.Mode.GetMode())
		s.True(s.tctx.mmuRegs.IF.VBlank.GetBool())
		s.Equal(uint8(0), s.tctx.mmuRegs.LY.Get())
	})
}

func (s *unitTestSuite) TestDoDataTransfer1() {
	s.Run("switch to data transfer 2", func() {
		s.tctx.mmuRegs.STAT.Mode.SetMode(registers.STATModeDataTransfer1)

		s.Require().NoError(s.tctx.gpu.DoCycles(0xFF))
		s.Equal(registers.STATModeDataTransfer2, s.tctx.mmuRegs.STAT.Mode.GetMode())
	})
}

func (s *unitTestSuite) TestDoDataTransfer2() {
	s.Run("switch to h-blank", func() {
		s.tctx.mmuRegs.STAT.Mode.SetMode(registers.STATModeDataTransfer2)

		s.Require().NoError(s.tctx.gpu.DoCycles(0xFF))
		s.Equal(registers.STATModeHBlank, s.tctx.mmuRegs.STAT.Mode.GetMode())
	})
}
