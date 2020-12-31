package memory

import (
	"errors"
)

func (s *unitTestSuite) TestWrapReadError_Nil() {
	s.NoError(wrapReadError(0x100, nil))
}

func (s *unitTestSuite) TestWrapReadError_NotNil() {
	s.EqualError(
		wrapReadError(0x100, errors.New("failed to read data")),
		"read error at 0100h: failed to read data",
	)
}

func (s *unitTestSuite) TestWrapWriteError_Nil() {
	s.NoError(wrapWriteError(0x100, nil))
}

func (s *unitTestSuite) TestWrapWriteError_NotNil() {
	s.EqualError(
		wrapWriteError(0x100, errors.New("failed to write data")),
		"write error at 0100h: failed to write data",
	)
}
