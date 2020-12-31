package registers

import (
	"nebula-go/pkg/gbc/memory/lib"
)

type DMAReg struct {
	mmu lib.MemoryIO
}

func NewDMAReg(mmu lib.MemoryIO) *DMAReg {
	return &DMAReg{
		mmu: mmu,
	}
}

func (r *DMAReg) Get() (uint8, error) {
	return 0, nil
}

func (r *DMAReg) Set(value uint8) error {
	if 0xF1 < value {
		value = 0xF1
	}

	startAddr := uint16(value) << 8

	values, err := r.mmu.ReadByteSlice(startAddr, 0xA0)
	if err != nil {
		return err
	}
	return r.mmu.WriteByteSlice(0xFE00, values)
}
