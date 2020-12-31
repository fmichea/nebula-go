package mbcs

import (
	"nebula-go/pkg/gbc/memory/segments"
)

// FIXME: Implementatin directly copied from C++ version without checking, need
//  to verify this is the expected behavior.

type mbc2 struct {
	defaultImpl
}

func NewMBC2(rom, eram segments.Segment) MBC {
	return newMBCWrapper(rom, eram, &mbc2{
		defaultImpl: defaultImpl{
			rom:  rom,
			eram: eram,
		},
	})
}

func (m *mbc2) bankSelectorZone1(addr uint16, value uint8) error {
	value &= 0x0F
	if value == 0 {
		value = 1
	}
	return m.rom.SelectBank(uint(value))
}

func (m *mbc2) bankSelectorZone2(value uint8) error {
	return nil
}

func (m *mbc2) bankModeSelect(value uint8) error {
	return nil
}
