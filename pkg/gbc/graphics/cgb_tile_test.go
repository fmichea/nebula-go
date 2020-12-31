package graphics

import (
	"nebula-go/pkg/common/frontends"
)

func (s *unitTestSuite) TestNewCGBBGWTile() {
	var (
		//data1 = make([]byte, 8*8*2)

		attrs1 = NewTileAttributes(0)
		attrs2 = NewTileAttributes(0x4)
	)

	s.Run("background tile uses BGPD palette", func() {
		tile := NewCGBBGWTile(s.tctx.mockMMU, 0, 0, attrs1)
		s.Require().NotNil(tile)

		startX, endX, pixels := tile.Colors(0)
		s.Equal(int16(0), startX)
		s.Equal(int16(8), endX)
		s.Len(pixels, 8)

		s.Equal(frontends.White, pixels[0].Color)
	})

	s.Run("the palette index in attrs is used", func() {
		tile := NewCGBBGWTile(s.tctx.mockMMU, 0, 0, attrs2)
		s.Require().NotNil(tile)

		s.tctx.mmuRegs.BGPI.AutoIncrement.SetBool(true)
		s.tctx.mmuRegs.BGPI.Index.Set(0x20)
		for x := 0; x < 12; x++ {
			s.tctx.mmuRegs.BGPD.Set(uint8(x + 1))
		}

		startX, endX, pixels := tile.Colors(0)
		s.Equal(int16(0), startX)
		s.Equal(int16(8), endX)
		s.Len(pixels, 8)

		s.NotEqual(frontends.Black, pixels[0].Color)
	})
}

func (s *unitTestSuite) TestNewCGBObjTile() {
	var (
		//data1 = make([]byte, 8*8*2)

		attrs1 = NewTileAttributes(0)
	)

	s.Run("object tile uses OBPD palette", func() {
		tile := NewCGBObjTile(s.tctx.mockMMU, 0, 0, 8, attrs1)
		s.Require().NotNil(tile)

		startX, endX, pixels := tile.Colors(0)
		s.Equal(int16(0), startX)
		s.Equal(int16(8), endX)
		s.Len(pixels, 8)

		s.Equal(frontends.White, pixels[0].Color)
	})
}
