package registers

import (
	"nebula-go/pkg/gbc/memory/cartridge"
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
	"nebula-go/pkg/gbc/memory/segments"
)

type VBKReg struct {
	registerslib.ByteWithErrors

	cr  *cartridge.ROM
	seg segments.Segment
}

func findMaskFromBankCount(bankCount uint) uint8 {
	maskPower := uint(0x2)
	for maskPower < bankCount {
		maskPower *= 2
	}
	return uint8(maskPower - 1)
}

func NewVBKReg(ptr *uint8, cr *cartridge.ROM, seg segments.Segment) (*VBKReg, error) {
	mask := findMaskFromBankCount(seg.BankCount())

	reg := registerslib.NewByteWithMask(ptr, 0x00, mask)

	vbk := &VBKReg{
		ByteWithErrors: registerslib.WrapWithError(reg),

		cr:  cr,
		seg: seg,
	}
	return vbk, vbk.Set(0)
}

func (r *VBKReg) Set(value uint8) error {
	if err := r.ByteWithErrors.Set(value); err != nil {
		return err
	}

	if r.cr.IsCGB() {
		value, err := r.Get()
		if err != nil {
			return err
		}

		if err := r.seg.SelectBank(uint(value)); err != nil {
			return err
		}
	}

	return nil
}
