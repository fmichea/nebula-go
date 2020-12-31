package graphics

import (
	"nebula-go/pkg/common/testhelpers"
	"nebula-go/pkg/gbc/memory/lib"

	"github.com/golang/mock/gomock"
)

// FIXME: document all hardcoded values, remove gomock.Any()

func buildTileData() []byte {
	data := make([]byte, 16*8*2)

	for idx := range data {
		data[idx] = byte(idx % 0xFF)
	}
	return data
}

var (
	_fillSprite = []byte("\xFF\x00\x00\x00")

	_validSprite1               = []byte("\x12\x02\x00\x00")
	_validSprite1Addr    uint16 = 0x8000
	_validSprite1Addr16H uint16 = 0x8000

	_validSprite2               = []byte("\x12\x10\x04\x00")
	_validSprite2Addr    uint16 = 0x8040
	_validSprite2Addr16H uint16 = 0x8080

	_validSprite3               = []byte("\x12\x50\x05\xFF")
	_validSprite3Addr16H uint16 = 0x8080

	_validSpriteHiddenX  = []byte("\x12\xFF\x09\x00")
	_validSpriteHiddenX2 = []byte("\x12\x00\x09\x00")

	_tileData = buildTileData()
)

func (s *unitTestSuite) buildSpriteFillSpace(count int) []byte {
	var result []byte

	for idx := 0; idx < count; idx++ {
		result = append(result, _fillSprite...)
	}
	return result
}

func (s *unitTestSuite) buildSpriteBuffer(buffers ...[]byte) []byte {
	var buffer []byte

	for _, b := range buffers {
		buffer = append(buffer, b...)
	}

	fillCount := (_spriteAttributesTableByteSize - len(buffer)) / _spriteAttributesByteSize
	if 0 < fillCount {
		buffer = append(buffer, s.buildSpriteFillSpace(fillCount)...)
	}

	s.Require().Len(buffer, _spriteAttributesTableByteSize)
	return buffer
}

func (s *unitTestSuite) TestLoadSpritesForLine() {
	var genericBuffer = s.buildSpriteBuffer(
		_validSprite2,
		s.buildSpriteFillSpace(5),
		_validSprite1,
		s.buildSpriteFillSpace(5),
		_validSpriteHiddenX,
	)

	s.Run("memory error is error", func() {
		s.tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xFE00), uint(0xA0)).Return(nil, testhelpers.ErrTesting1)

		sprites, err := s.tctx.gpu.loadSpritesForLine(0)
		s.Error(err)
		s.Empty(sprites)
	})

	s.Run("no matching sprite is empty result", func() {
		s.tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xFE00), uint(0xA0)).Return(genericBuffer, nil)

		sprites, err := s.tctx.gpu.loadSpritesForLine(0)
		s.NoError(err)
		s.Empty(sprites)
	})

	s.Run("two sprites on line are returned, one hidden X ignored", func() {
		s.tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xFE00), uint(0xA0)).Return(genericBuffer, nil)

		s.tctx.mockMMU.EXPECT().ReadDByte(gomock.Any()).Return(uint16(0), nil)
		s.tctx.mockMMU.EXPECT().ReadDByte(gomock.Any()).Return(uint16(0), nil)

		sprites, err := s.tctx.gpu.loadSpritesForLine(4)
		s.NoError(err)
		s.Len(sprites, 2)
	})

	s.Run("hidden X affects how many sprites can be fecthed, but not fills", func() {
		buffer := s.buildSpriteBuffer(
			_validSpriteHiddenX,
			_validSpriteHiddenX2,
			_validSpriteHiddenX,
			s.buildSpriteFillSpace(5),
			_validSpriteHiddenX,
			_validSprite1,
			_validSpriteHiddenX2,
			_validSpriteHiddenX,
			s.buildSpriteFillSpace(5),
			_validSpriteHiddenX2,
			_validSpriteHiddenX2,
			_validSpriteHiddenX,
			_validSprite2,
		)

		s.tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xFE00), uint(0xA0)).Return(buffer, nil)

		s.tctx.mockMMU.EXPECT().ReadDByte(gomock.Any()).Return(uint16(0), nil)

		sprites, err := s.tctx.gpu.loadSpritesForLine(4)
		s.NoError(err)
		s.Len(sprites, 1)
	})

	s.Run("y limit for 8 pixels high sprite", func() {
		s.tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xFE00), uint(0xA0)).Return(genericBuffer, nil).Times(2)
		s.tctx.mockMMU.EXPECT().ReadDByte(gomock.Any()).Return(uint16(0), nil)
		s.tctx.mockMMU.EXPECT().ReadDByte(gomock.Any()).Return(uint16(0), nil)

		sprites, err := s.tctx.gpu.loadSpritesForLine(9)
		s.NoError(err)
		s.Len(sprites, 2)

		sprites, err = s.tctx.gpu.loadSpritesForLine(10)
		s.NoError(err)
		s.Len(sprites, 0) // no sprites because we are past the 8 pixels high sprites.
	})

	s.Run("y limit for 16 pixels high sprite", func() {
		s.tctx.mmuRegs.LCDC.OBJSS.SetBool(true)

		s.tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xFE00), uint(0xA0)).Return(genericBuffer, nil).Times(2)
		s.tctx.mockMMU.EXPECT().ReadDByte(gomock.Any()).Return(uint16(0), nil).Times(2)
		s.tctx.mockMMU.EXPECT().ReadDByte(gomock.Any()).Return(uint16(0), nil).Times(2)

		sprites, err := s.tctx.gpu.loadSpritesForLine(9)
		s.NoError(err)
		s.Len(sprites, 2)

		sprites, err = s.tctx.gpu.loadSpritesForLine(10)
		s.NoError(err)
		s.Len(sprites, 2) // here we get two sprites, because they are 16 pixels high.
	})

	s.Run("invalid sprite data read is error", func() {
		s.tctx.mmuRegs.LCDC.OBJSS.SetBool(true)

		s.tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xFE00), uint(0xA0)).Return(genericBuffer, nil)
		s.tctx.mockMMU.EXPECT().ReadDByte(gomock.Any()).Return(uint16(0), testhelpers.ErrTesting1)

		_, err := s.tctx.gpu.loadSpritesForLine(9)
		s.Equal(testhelpers.ErrTesting1, err)
	})

	s.Run("miss-aligned sprite with 16 lines height is re-aligned", func() {
		buffer := s.buildSpriteBuffer(_validSprite3)

		s.tctx.mmuRegs.LCDC.OBJSS.SetBool(true)
		s.tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xFE00), uint(0xA0)).Return(buffer, nil).Times(2)

		//var data []byte
		//for idx := 0; idx < 8; idx++ {
		//	data = append(data, 0xFF, 0x0F)
		//}
		//for idx := 0; idx < 8; idx++ {
		//	data = append(data, 0xF0, 0xFF)
		//}

		s.tctx.mockMMU.EXPECT().ReadDByte(gomock.Any()).Return(uint16(0xFFF0), nil)

		sprites, err := s.tctx.gpu.loadSpritesForLine(4)
		s.NoError(err)
		s.Require().Len(sprites, 1)

		var colorIDs []uint8

		_, _, pixels := sprites[0].Colors(4)
		for _, pixel := range pixels {
			colorIDs = append(colorIDs, pixel.ColorID)
		}
		s.Equal([]uint8{0x2, 0x2, 0x2, 0x2, 0x3, 0x3, 0x3, 0x3}, colorIDs)

		s.tctx.mockMMU.EXPECT().ReadDByte(gomock.Any()).Return(uint16(0x0FFF), nil)

		sprites, err = s.tctx.gpu.loadSpritesForLine(12)
		s.NoError(err)
		s.Require().Len(sprites, 1)

		colorIDs = []uint8{}
		_, _, pixels = sprites[0].Colors(12)
		for _, pixel := range pixels {
			colorIDs = append(colorIDs, pixel.ColorID)
		}
		s.Equal([]uint8{0x3, 0x3, 0x3, 0x3, 0x1, 0x1, 0x1, 0x1}, colorIDs)
	})

	s.Run("DMG sorts tiles by X", func() {
		s.tctx.cr.Type = lib.DMG01

		s.tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xFE00), uint(0xA0)).Return(genericBuffer, nil)

		s.tctx.mockMMU.EXPECT().ReadDByte(gomock.Any()).Return(uint16(0), nil)
		s.tctx.mockMMU.EXPECT().ReadDByte(gomock.Any()).Return(uint16(0), nil)

		sprites, err := s.tctx.gpu.loadSpritesForLine(4)
		s.Require().NoError(err)
		s.Require().Len(sprites, 2)

		startX1, _, _ := sprites[0].Colors(4)
		startX2, _, _ := sprites[1].Colors(4)
		s.Less(startX1, startX2)
	})

	s.Run("CGB keeps the sprites in order found", func() {
		s.tctx.cr.Type = lib.CGB001

		s.tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xFE00), uint(0xA0)).Return(genericBuffer, nil)

		s.tctx.mockMMU.EXPECT().ReadDByte(gomock.Any()).Return(uint16(0), nil)
		s.tctx.mockMMU.EXPECT().ReadDByte(gomock.Any()).Return(uint16(0), nil)

		sprites, err := s.tctx.gpu.loadSpritesForLine(4)
		s.Require().NoError(err)
		s.Require().Len(sprites, 2)

		startX1, _, _ := sprites[0].Colors(4)
		startX2, _, _ := sprites[1].Colors(4)
		s.Less(startX2, startX1)
	})
}
