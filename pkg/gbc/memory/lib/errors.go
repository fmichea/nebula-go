package lib

import (
	"errors"
)

var (
	ErrMBCNotImplemented = errors.New("memory controller has not been implemented yet")
	ErrInvalidRead       = errors.New("invalid read")
	ErrInvalidWrite      = errors.New("invalid write")
)
