package mbcs

import (
	"nebula-go/pkg/gbc/memory/lib"
)

func (s *unitTestSuite) TestROMOnly_ContainsAddress() {
	s.testContainsAddress(s.romOnly)
}

func (s *unitTestSuite) TestROMOnly_ROMReadFunctional() {
	ptr, err := s.romOnly.BytePtr(lib.AccessTypeRead, 0x4000, 0)
	s.NoError(err)
	s.NotNil(ptr)
	s.Equal(uint8(0x01), *ptr)
}

func (s *unitTestSuite) TestROMOnly_CannotReadRAM() {
	ptr, err := s.romOnly.BytePtr(lib.AccessTypeRead, 0xA000, 0)
	s.Equal(lib.ErrInvalidRead, err)
	s.Nil(ptr)
}

func (s *unitTestSuite) TestROMOnly_CannotWriteROM() {
	ptr, err := s.romOnly.BytePtr(lib.AccessTypeWrite, 0x4000, 0)
	s.Equal(lib.ErrInvalidWrite, err)
	s.Nil(ptr)
}

func (s *unitTestSuite) TestROMOnly_CannotWriteRAM() {
	ptr, err := s.romOnly.BytePtr(lib.AccessTypeWrite, 0xA000, 0)
	s.Equal(lib.ErrInvalidWrite, err)
	s.Nil(ptr)
}
