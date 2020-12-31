package graphics

import (
	"nebula-go/pkg/common/frontends"
)

func (s *unitTestSuite) TestDMGTile() {
	var (
		tileAddr uint16 = 0xABCD
		tileData uint16 = 0xFFFF
	)

	s.Run("BGW tile uses the BGP palette", func() {
		s.tctx.mmuRegs.BGP.Set(0xFF)

		s.tctx.mmuRegs.OBP0.Set(0x00)
		s.tctx.mmuRegs.OBP1.Set(0xAA)

		s.tctx.mockMMU.EXPECT().ReadDByte(tileAddr).Return(tileData, nil)

		tile := NewDMGBGWTile(s.tctx.mockMMU, 0, 0)
		s.Require().NoError(tile.LoadLineData(s.tctx.mockMMU, tileAddr, 0, 0))

		_, _, pixels := tile.Colors(0)

		s.Len(pixels, 8)
		s.Equal(uint8(0x3), pixels[0].ColorID)
		s.Equal(frontends.Black, pixels[0].Color)
	})

	s.Run("Object tile uses the first palette depending on attr: v0", func() {
		s.tctx.mmuRegs.BGP.Set(0xFF)

		s.tctx.mmuRegs.OBP0.Set(0x00)
		s.tctx.mmuRegs.OBP1.Set(0xAA)

		s.tctx.mockMMU.EXPECT().ReadDByte(tileAddr).Return(tileData, nil)

		attrs := NewTileAttributes(0)

		tile := NewDMGObjTile(s.tctx.mockMMU, 0, 0, 8, attrs)
		s.Require().NoError(tile.LoadLineData(s.tctx.mockMMU, tileAddr, 0, 0))
		_, _, pixels := tile.Colors(0)

		s.Len(pixels, 8)
		s.Equal(uint8(0x3), pixels[0].ColorID)
		s.Equal(frontends.White, pixels[0].Color)
	})

	s.Run("Object tile uses the first palette depending on attr: v1", func() {
		s.tctx.mmuRegs.BGP.Set(0xFF)

		s.tctx.mmuRegs.OBP0.Set(0x00)
		s.tctx.mmuRegs.OBP1.Set(0xAA)

		s.tctx.mockMMU.EXPECT().ReadDByte(tileAddr).Return(tileData, nil)

		attrs := NewTileAttributes(0x10)

		tile := NewDMGObjTile(s.tctx.mockMMU, 0, 0, 8, attrs)
		s.Require().NoError(tile.LoadLineData(s.tctx.mockMMU, tileAddr, 0, 0))
		_, _, pixels := tile.Colors(0)

		s.Len(pixels, 8)
		s.Equal(uint8(0x3), pixels[0].ColorID)
		s.Equal(frontends.NewPixel(0x55, 0x55, 0x55), pixels[0].Color)
	})
}
