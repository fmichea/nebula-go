package mbcs

import (
	"nebula-go/pkg/gbc/memory/lib"
)

func (s *unitTestSuite) TestMBCWrapper_InvalidRead() {
	ptr, err := s.mbc1.BytePtr(lib.AccessTypeRead, 0xFFFF, 0)
	s.Equal(lib.ErrInvalidRead, err)
	s.Nil(ptr)
}

func (s *unitTestSuite) TestMBCWrapper_InvalidWrite() {
	ptr, err := s.mbc1.BytePtr(lib.AccessTypeWrite, 0xFFFF, 0)
	s.Equal(lib.ErrInvalidWrite, err)
	s.Nil(ptr)
}
