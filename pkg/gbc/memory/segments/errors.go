package segments

import (
	"errors"
)

var (
	ErrBankUnavailable       = errors.New("bank selected is not available")
	ErrBankCountInvalid      = errors.New("invalid number of banks in configuration")
	ErrInvalidMirrorRange    = errors.New("mirrored range does not fit in segment")
	ErrBufferIncompatible    = errors.New("raw buffer was incompatible with internal segment buffer")
	ErrCannotPin0WithOneBank = errors.New("cannot have a pinned bank 0 and only 1 bank")
	ErrSegmentTooSmall       = errors.New("cannot read/write slice: segment too small")
	ErrInvalidHookInBank     = errors.New("invalid attempt at byte hook in banked memory")
)
