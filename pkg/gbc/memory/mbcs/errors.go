package mbcs

import (
	"errors"
)

var (
	ErrMBCHookInvalid          = errors.New("MBC cannot be hooked")
	ErrRAMUnavailable          = errors.New("RAM is not available at the moment")
	ErrMBCSliceOperatioInvalid = errors.New("MBC does not support slice operations")
)
