package graphics

import (
	"github.com/golang/mock/gomock"

	"nebula-go/pkg/common/frontends"
	"nebula-go/pkg/common/testhelpers"
	graphicslib "nebula-go/pkg/gbc/graphics/lib"
	"nebula-go/pkg/gbc/memory/lib"
)

func (s *unitTestSuite) TestShouldRenderBackground() {
	s.Run("CGB always renders background", func() {
		s.tctx.cr.Type = lib.CGB001

		s.tctx.mmuRegs.LCDC.BGD.SetBool(true)
		s.True(s.tctx.gpu.shouldRenderBackground())

		s.tctx.mmuRegs.LCDC.BGD.SetBool(false)
		s.True(s.tctx.gpu.shouldRenderBackground())
	})

	s.Run("DMG renders background only when LCDC bit 0 is set", func() {
		s.tctx.cr.Type = lib.DMG01

		s.tctx.mmuRegs.LCDC.BGD.SetBool(true)
		s.True(s.tctx.gpu.shouldRenderBackground())

		s.tctx.mmuRegs.LCDC.BGD.SetBool(false)
		s.False(s.tctx.gpu.shouldRenderBackground())
	})
}

func (s *unitTestSuite) TestGetBackgroundPixels() {
	var (
		tile1ID   uint8  = 0
		tile1Data uint16 = 0
	)

	s.Run("DMG and disabled returns an empty background", func() {
		s.tctx.cr.Type = lib.DMG01
		s.tctx.mmuRegs.LCDC.BGD.SetBool(false)

		pixels, err := s.tctx.gpu.getBackgroundPixels(0)
		s.Require().NoError(err)
		s.Len(pixels, int(graphicslib.Width))
		s.Nil(pixels[0])
	})

	s.Run("render loads up all of the tiles and gets relevant pixels", func() {
		s.tctx.cr.Type = lib.DMG01
		s.tctx.mmuRegs.LCDC.BGD.SetBool(true)

		for idx := 0; idx < 20; idx++ {
			s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9800+idx)).Return(tile1ID, nil)
		}
		s.tctx.mockMMU.EXPECT().ReadDByte(uint16(0x8000)).Return(tile1Data, nil).Times(20)

		pixels, err := s.tctx.gpu.getBackgroundPixels(0)
		s.Require().NoError(err)
		s.Len(pixels, int(graphicslib.Width))
		s.NotNil(pixels[0])
	})

	s.Run("CGB render loads up all of the tiles and gets relevant pixels even disabled", func() {
		s.tctx.cr.Type = lib.CGB001
		s.tctx.mmuRegs.LCDC.BGD.SetBool(false)

		s.tctx.mockVRAM.EXPECT().SelectBank(gomock.Any()).Return(nil).AnyTimes()
		for idx := 0; idx < 20; idx++ {
			s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9800+idx)).Return(tile1ID, nil).Times(2)
		}
		s.tctx.mockMMU.EXPECT().ReadDByte(uint16(0x8000)).Return(tile1Data, nil).Times(20)

		pixels, err := s.tctx.gpu.getBackgroundPixels(0)
		s.Require().NoError(err)
		s.Len(pixels, int(graphicslib.Width))
		s.NotNil(pixels[0])
	})

	s.Run("background tile failure to load", func() {
		s.tctx.cr.Type = lib.DMG01
		s.tctx.mmuRegs.LCDC.BGD.SetBool(true)

		s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9800)).Return(uint8(0), testhelpers.ErrTesting1)

		_, err := s.tctx.gpu.getBackgroundPixels(0)
		s.Equal(testhelpers.ErrTesting1, err)
	})

	// FIXME: test wrap around here.
}

func (s *unitTestSuite) TestGetWindowPixels() {
	var (
		tile1ID   uint8  = 0
		tile1Data uint16 = 0
	)

	s.Run("enabled loads tile for the window", func() {
		s.tctx.cr.Type = lib.DMG01
		s.tctx.mmuRegs.LCDC.WDE.SetBool(true)

		s.tctx.mmuRegs.WX.Set(17)
		s.tctx.mmuRegs.WY.Set(7)

		for idx := 1; idx < 20; idx++ {
			s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9800+idx)).Return(tile1ID, nil)
		}
		s.tctx.mockMMU.EXPECT().ReadDByte(uint16(0x8006)).Return(tile1Data, nil).Times(19)

		pixels, err := s.tctx.gpu.getWindowPixels(10)
		s.Require().NoError(err)
		s.Len(pixels, int(graphicslib.Width))
		s.Nil(pixels[0])
		s.NotNil(pixels[10])
		s.NotNil(pixels[graphicslib.Width-1])
	})

	s.Run("disabled returns empty pixels", func() {
		s.tctx.cr.Type = lib.DMG01
		s.tctx.mmuRegs.LCDC.WDE.SetBool(false)

		s.tctx.mmuRegs.WX.Set(17)
		s.tctx.mmuRegs.WY.Set(7)

		pixels, err := s.tctx.gpu.getWindowPixels(10)
		s.Require().NoError(err)
		s.Len(pixels, int(graphicslib.Width))
		s.Nil(pixels[0])
		s.Nil(pixels[10])
		s.Nil(pixels[graphicslib.Width-1])
	})

	s.Run("rendering line before the window returns empty", func() {
		s.tctx.mmuRegs.LCDC.WDE.SetBool(true)

		s.tctx.mmuRegs.WX.Set(17)
		s.tctx.mmuRegs.WY.Set(7)

		pixels, err := s.tctx.gpu.getWindowPixels(0)
		s.Require().NoError(err)
		s.Len(pixels, int(graphicslib.Width))
		s.Nil(pixels[0])
		s.Nil(pixels[10])
		s.Nil(pixels[graphicslib.Width-1])
	})
}

func (s *unitTestSuite) TestGetObjectPixels() {
	s.Run("objects disabled returns empty pixels", func() {
		s.tctx.mmuRegs.LCDC.OBJSDE.SetBool(false)

		pixels, err := s.tctx.gpu.getObjectPixels(0)
		s.Require().NoError(err)
		s.Len(pixels, int(graphicslib.Width))
		s.Nil(pixels[0])
		s.Nil(pixels[10])
		s.Nil(pixels[graphicslib.Width-1])
	})

	s.Run("objects are drawn only where available", func() {
		s.tctx.mmuRegs.LCDC.OBJSDE.SetBool(true)

		buffer := s.buildSpriteBuffer(_validSprite1)
		s.tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xFE00), uint(0xA0)).Return(buffer, nil)

		s.tctx.mockMMU.EXPECT().ReadDByte(_validSprite1Addr+0x4).Return(uint16(0x0002), nil)

		pixels, err := s.tctx.gpu.getObjectPixels(4)
		s.Require().NoError(err)

		s.NotNil(pixels[0]) // pixel 7 of tile, at pixel 0 of line, color id 2, not transparent.
		s.Nil(pixels[1])    // pixel 8 of tile, at pixel 1 of line, color id 0, transparent.
	})
}

func (s *unitTestSuite) TestSelectPixelColor() {
	var (
		c1 = frontends.NewPixel(0xFF, 0x00, 0x00)
		c2 = frontends.NewPixel(0x00, 0xFF, 0x00)
		c3 = frontends.NewPixel(0x00, 0x00, 0xFF)

		p1 = &Pixel{
			BackgroundPriority: false,
			ColorID:            1,
			Color:              c1,
		}

		p2 = &Pixel{
			BackgroundPriority: false,
			ColorID:            2,
			Color:              c2,
		}

		p3 = &Pixel{
			BackgroundPriority: false,
			ColorID:            3,
			Color:              c3,
		}
	)

	s.Run("no pixels is white", func() {
		pixel := s.tctx.gpu.selectPixelColor(nil, nil, nil)
		s.Equal(frontends.White, pixel)
	})

	s.Run("only background available is returned", func() {
		pixel := s.tctx.gpu.selectPixelColor(p1, nil, nil)
		s.Equal(c1, pixel)
	})

	s.Run("only window available is returned", func() {
		pixel := s.tctx.gpu.selectPixelColor(nil, p1, nil)
		s.Equal(c1, pixel)
	})

	s.Run("only obj pixel available is returned", func() {
		pixel := s.tctx.gpu.selectPixelColor(nil, nil, p1)
		s.Equal(c1, pixel)
	})

	s.Run("background is below window, no object", func() {
		pixel := s.tctx.gpu.selectPixelColor(p1, p2, nil)
		s.Equal(c2, pixel)
	})

	s.Run("background is below obj, no window", func() {
		pixel := s.tctx.gpu.selectPixelColor(p1, nil, p3)
		s.Equal(c3, pixel)
	})

	s.Run("window is below obj, no background", func() {
		pixel := s.tctx.gpu.selectPixelColor(nil, p2, p3)
		s.Equal(c3, pixel)
	})

	s.Run("CGB: LCDC.BG Display enabled and BG priority on BG pixel is override", func() {
		s.tctx.mmuRegs.LCDC.BGD.SetBool(true)

		p1 = &Pixel{
			BackgroundPriority: true,
			ColorID:            1,
			Color:              c1,
		}

		pixel := s.tctx.gpu.selectPixelColor(p1, nil, p3)
		s.Equal(c1, pixel)
	})

	s.Run("CGB: LCDC.BG Display enabled and window priority is override", func() {
		s.tctx.mmuRegs.LCDC.BGD.SetBool(true)

		p2 := &Pixel{
			BackgroundPriority: true,
			ColorID:            2,
			Color:              c2,
		}

		pixel := s.tctx.gpu.selectPixelColor(nil, p2, p3)
		s.Equal(c2, pixel)
	})

	s.Run("CGB: LCDC.BG Display enabled and BG priority but window on top, priority to obj", func() {
		s.tctx.mmuRegs.LCDC.BGD.SetBool(true)

		p1 = &Pixel{
			BackgroundPriority: true,
			ColorID:            1,
			Color:              c1,
		}

		pixel := s.tctx.gpu.selectPixelColor(p1, p2, p3)
		s.Equal(c3, pixel)
	})

	s.Run("CGB: LCDC.BG Display disabled removes BG priority", func() {
		s.tctx.mmuRegs.LCDC.BGD.SetBool(false)

		p1 = &Pixel{
			BackgroundPriority: true,
			ColorID:            1,
			Color:              c1,
		}

		pixel := s.tctx.gpu.selectPixelColor(p1, nil, p3)
		s.Equal(c3, pixel)
	})

	s.Run("object with low priority with background color id 0", func() {
		p1 := &Pixel{
			BackgroundPriority: false,
			ColorID:            0,
			Color:              c1,
		}

		p3 := &Pixel{
			BackgroundPriority: true,
			ColorID:            2,
			Color:              c3,
		}

		pixel := s.tctx.gpu.selectPixelColor(p1, nil, p3)
		s.Equal(c3, pixel)
	})

	s.Run("object with low priority with background color id 1+", func() {
		p1 := &Pixel{
			BackgroundPriority: false,
			ColorID:            1,
			Color:              c1,
		}

		p3 := &Pixel{
			BackgroundPriority: true,
			ColorID:            2,
			Color:              c3,
		}

		pixel := s.tctx.gpu.selectPixelColor(p1, nil, p3)
		s.Equal(c1, pixel)
	})
}

func (s *unitTestSuite) TestDrawCurrentLine() {
	s.tctx.cr.Type = lib.DMG01

	s.tctx.mmuRegs.LY.Set(0)
	s.tctx.mmuRegs.LCDC.BGD.SetBool(false)
	s.tctx.mmuRegs.LCDC.WDE.SetBool(false)
	s.tctx.mmuRegs.LCDC.OBJSDE.SetBool(false)

	s.tctx.mockWindow.EXPECT().DrawLine(uint(0), gomock.Any()).Return(nil)

	s.NoError(s.tctx.gpu.drawCurrentLine())
}
