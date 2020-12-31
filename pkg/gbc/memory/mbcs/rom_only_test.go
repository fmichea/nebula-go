package mbcs

import (
	"nebula-go/pkg/gbc/memory/lib"
)

func (s *unitTestSuite) TestROMOnly_ContainsAddress() {
	s.testContainsAddress(s.romOnly)
}

func (s *unitTestSuite) TestROMOnly_ROMReadFunctional() {
	value, err := s.romOnly.ReadByte(0x4000)
	s.Require().NoError(err)
	s.Equal(uint8(0x01), value)
}

func (s *unitTestSuite) TestROMOnly_CannotReadRAM() {
	_, err := s.romOnly.ReadByte(0xA000)
	s.Require().Equal(lib.ErrInvalidRead, err)
}

func (s *unitTestSuite) TestROMOnly_CannotWriteROM() {
	err := s.romOnly.WriteByte(0x4000, 0)
	s.Require().Equal(lib.ErrInvalidWrite, err)
}

func (s *unitTestSuite) TestROMOnly_CannotWriteRAM() {
	err := s.romOnly.WriteByte(0xA000, 0)
	s.Require().Equal(lib.ErrInvalidWrite, err)
}

func (s *unitTestSuite) TestROMOnly_ReadByteSlice_OKInROM() {
	_, err := s.romOnly.ReadByteSlice(0x1000, 0x10)
	s.Require().NoError(err)
}

func (s *unitTestSuite) TestROMOnly_ReadByteSlice_InvalidInRAM() {
	_, err := s.romOnly.ReadByteSlice(0xA000, 0x10)
	s.Require().Equal(lib.ErrInvalidRead, err)
}

func (s *unitTestSuite) TestROMOnly_WriteByteSlice() {
	err := s.romOnly.WriteByteSlice(0x1000, nil)
	s.Require().Equal(ErrMBCSliceOperatioInvalid, err)
}

func (s *unitTestSuite) TestROMOnly_ByteHook() {
	_, err := s.romOnly.ByteHook(0x1000)
	s.Require().Equal(ErrMBCHookInvalid, err)
}
