package mbcs

import (
	"nebula-go/pkg/gbc/memory/lib"
	"nebula-go/pkg/gbc/memory/segments"
)

func (s *unitTestSuite) TestMBCWrapper_AddressRanges() {
	ranges := s.mbc1.AddressRanges()
	s.Len(ranges, 2)
	s.Equal(segments.AddressRange{Start: 0x0000, End: 0x7FFF}, ranges[0])
	s.Equal(segments.AddressRange{Start: 0xA000, End: 0xBFFF}, ranges[1])
}

func (s *unitTestSuite) TestMBCWrapper_InvalidRead() {
	_, err := s.mbc1.ReadByte(0xFFFF)
	s.Require().Equal(lib.ErrInvalidRead, err)
}

func (s *unitTestSuite) TestMBCWrapper_ValidInROM() {
	_, err := s.mbc1.ReadByteSlice(0x1000, 0x10)
	s.Require().NoError(err)
}

func (s *unitTestSuite) TestMBCWrapper_InvalidReadRAMDisabled() {
	_, err := s.mbc1.ReadByteSlice(0xA000, 0x10)
	s.Require().Equal(ErrRAMUnavailable, err)
}

func (s *unitTestSuite) TestMBCWrapper_InvalidWrite() {
	err := s.mbc1.WriteByte(0xFFFF, 0)
	s.Require().Equal(lib.ErrInvalidWrite, err)

	err = s.mbc1.WriteByteSlice(0x1000, nil)
	s.Require().Equal(ErrMBCSliceOperatioInvalid, err)
}

func (s *unitTestSuite) TestMBCWrapper_ByteHook() {
	_, err := s.mbc1.ByteHook(0x1000)
	s.Require().Equal(ErrMBCHookInvalid, err)
}
