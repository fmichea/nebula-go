package lib

import (
	"errors"
)

var (
	ErrMBCNotImplemented  = errors.New("memory controller has not been implemented yet")
	ErrInvalidRead        = errors.New("invalid read")
	ErrInvalidWrite       = errors.New("invalid write")
	ErrNoSegmentAtAddr    = errors.New("no segment at address")
	ErrInvalidSegmentAddr = errors.New("requested invalid address in segment")
	ErrDoubleHook         = errors.New("tried to hook address twice")
	ErrHookNotProvided    = errors.New("hook was not provided")
)
