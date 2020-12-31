package graphics

import (
	"github.com/golang/mock/gomock"

	"nebula-go/pkg/common/testhelpers"
	"nebula-go/pkg/gbc/memory/lib"
)

func (s *unitTestSuite) TestNewBackgroundTile() {
	//data := make([]byte, 16)

	s.Run("DMG tile is loaded fully", func() {
		s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(42), nil)
		s.tctx.mockMMU.EXPECT().ReadDByte(uint16(0x82a4)).Return(uint16(0), nil)

		tile, err := NewBackgroundTile(s.tctx.mockMMU, 40, 18)
		s.Require().NoError(err)
		s.NotNil(tile)
	})

	s.Run("CGB tile is loaded, with attributes from bank 1", func() {
		s.tctx.cr.Type = lib.CGB001
		s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(42), nil)

		c1 := s.tctx.mockVRAM.EXPECT().SelectBank(uint(1)).Return(nil)
		c2 := s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(0), nil)
		c3 := s.tctx.mockVRAM.EXPECT().SelectBank(uint(0)).Return(nil)
		gomock.InOrder(c1, c2, c3)

		s.tctx.mockMMU.EXPECT().ReadDByte(uint16(0x82a4)).Return(uint16(0), nil)

		tile, err := NewBackgroundTile(s.tctx.mockMMU, 40, 18)
		s.Require().NoError(err)
		s.NotNil(tile)
	})

	s.Run("CGB tile is loaded, with attributes, changes bank for data", func() {
		s.tctx.cr.Type = lib.CGB001
		s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(42), nil)

		c1 := s.tctx.mockVRAM.EXPECT().SelectBank(uint(1)).Return(nil)
		c2 := s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(0x8), nil)
		c3 := s.tctx.mockVRAM.EXPECT().SelectBank(uint(0)).Return(nil)
		gomock.InOrder(c1, c2, c3)

		c1 = s.tctx.mockVRAM.EXPECT().SelectBank(uint(1)).Return(nil)
		c2 = s.tctx.mockMMU.EXPECT().ReadDByte(uint16(0x82a4)).Return(uint16(0), nil)
		c3 = s.tctx.mockVRAM.EXPECT().SelectBank(uint(0)).Return(nil)
		gomock.InOrder(c1, c2, c3)

		tile, err := NewBackgroundTile(s.tctx.mockMMU, 40, 18)
		s.Require().NoError(err)
		s.NotNil(tile)
	})

	s.Run("tmds changes where the tile number is loaded from", func() {
		s.tctx.mmuRegs.LCDC.BGTMDS.Set(1)
		s.tctx.mmuRegs.LCDC.WTMDS.Set(0)

		s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9C45)).Return(uint8(42), nil)
		s.tctx.mockMMU.EXPECT().ReadDByte(uint16(0x82a4)).Return(uint16(0), nil)

		tile, err := NewBackgroundTile(s.tctx.mockMMU, 40, 18)
		s.Require().NoError(err)
		s.NotNil(tile)
	})

	s.Run("DMG get tile identifier failure is overall error", func() {
		s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(0), testhelpers.ErrTesting1)

		_, err := NewBackgroundTile(s.tctx.mockMMU, 40, 18)
		s.Equal(testhelpers.ErrTesting1, err)
	})

	s.Run("DMG load tile data error is overall error", func() {
		s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(42), nil)
		s.tctx.mockMMU.EXPECT().ReadDByte(uint16(0x82a4)).Return(uint16(0), testhelpers.ErrTesting1)

		_, err := NewBackgroundTile(s.tctx.mockMMU, 40, 18)
		s.Equal(testhelpers.ErrTesting1, err)
	})

	s.Run("CGB load tile identifier error is overall error", func() {
		s.tctx.cr.Type = lib.CGB001
		s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(0), testhelpers.ErrTesting1)

		_, err := NewBackgroundTile(s.tctx.mockMMU, 40, 18)
		s.Equal(testhelpers.ErrTesting1, err)
	})

	s.Run("CGB vram bank change failure is overall error", func() {
		s.tctx.cr.Type = lib.CGB001
		s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(42), nil)

		s.tctx.mockVRAM.EXPECT().SelectBank(uint(1)).Return(testhelpers.ErrTesting1)

		_, err := NewBackgroundTile(s.tctx.mockMMU, 40, 18)
		s.Equal(testhelpers.ErrTesting1, err)
	})

	s.Run("CGB attributes read failure is overall error", func() {
		s.tctx.cr.Type = lib.CGB001
		s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(42), nil)

		s.tctx.mockVRAM.EXPECT().SelectBank(uint(1)).Return(nil)
		s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(0), testhelpers.ErrTesting1)

		_, err := NewBackgroundTile(s.tctx.mockMMU, 40, 18)
		s.Equal(testhelpers.ErrTesting1, err)
	})

	s.Run("CGB vram bank comeback error is failure", func() {
		s.tctx.cr.Type = lib.CGB001
		s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(42), nil)

		c1 := s.tctx.mockVRAM.EXPECT().SelectBank(uint(1)).Return(nil)
		c2 := s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(0x8), nil)
		c3 := s.tctx.mockVRAM.EXPECT().SelectBank(uint(0)).Return(testhelpers.ErrTesting1)
		gomock.InOrder(c1, c2, c3)

		_, err := NewBackgroundTile(s.tctx.mockMMU, 40, 18)
		s.Equal(testhelpers.ErrTesting1, err)
	})

	s.Run("CGB tile data error is overall error", func() {
		s.tctx.cr.Type = lib.CGB001
		s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(42), nil)

		c1 := s.tctx.mockVRAM.EXPECT().SelectBank(uint(1)).Return(nil)
		c2 := s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(0x8), nil)
		c3 := s.tctx.mockVRAM.EXPECT().SelectBank(uint(0)).Return(nil)
		gomock.InOrder(c1, c2, c3)

		c1 = s.tctx.mockVRAM.EXPECT().SelectBank(uint(1)).Return(nil)
		c2 = s.tctx.mockMMU.EXPECT().ReadDByte(uint16(0x82a4)).Return(uint16(0), testhelpers.ErrTesting1)
		gomock.InOrder(c1, c2)

		_, err := NewBackgroundTile(s.tctx.mockMMU, 40, 18)
		s.Equal(testhelpers.ErrTesting1, err)
	})

	s.Run("x-axis invalid is denied", func() {
		_, err := NewBackgroundTile(s.tctx.mockMMU, 1000, 0)
		s.Equal(ErrBGWXCoordinateInvalid, err)
	})

	s.Run("y-axis invalid is denied", func() {
		_, err := NewBackgroundTile(s.tctx.mockMMU, 0, 1000)
		s.Equal(ErrBGWYCoordinateInvalid, err)
	})
}

func (s *unitTestSuite) TestNewWindowTile() {
	// NOTE: Background and window work the same, so this just tests that
	// the window only takes care of its own data select.

	//data := make([]byte, 16)

	s.Run("tmds changes where the tile number is loaded from", func() {
		s.tctx.mmuRegs.LCDC.BGTMDS.Set(1)
		s.tctx.mmuRegs.LCDC.WTMDS.Set(0)

		s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9845)).Return(uint8(42), nil)
		s.tctx.mockMMU.EXPECT().ReadDByte(uint16(0x82a4)).Return(uint16(0), nil)

		tile, err := NewWindowTile(s.tctx.mockMMU, 40, 18)
		s.Require().NoError(err)
		s.NotNil(tile)
	})

	s.Run("tmds changes where the tile number is loaded from", func() {
		s.tctx.mmuRegs.LCDC.BGTMDS.Set(0)
		s.tctx.mmuRegs.LCDC.WTMDS.Set(1)

		s.tctx.mockMMU.EXPECT().ReadByte(uint16(0x9C45)).Return(uint8(42), nil)
		s.tctx.mockMMU.EXPECT().ReadDByte(uint16(0x82a4)).Return(uint16(0), nil)

		tile, err := NewWindowTile(s.tctx.mockMMU, 40, 18)
		s.Require().NoError(err)
		s.NotNil(tile)
	})
}
