package memory

import (
	"fmt"
)

func wrapMMUError(t string, addr uint16, err error) error {
	if err != nil {
		return fmt.Errorf("%s error at %04Xh: %s", t, addr, err.Error())
	}
	return nil
}

func wrapReadError(addr uint16, err error) error {
	return wrapMMUError("read", addr, err)
}

func wrapWriteError(addr uint16, err error) error {
	return wrapMMUError("write", addr, err)
}
