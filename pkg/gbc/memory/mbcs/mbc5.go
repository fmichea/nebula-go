package mbcs

import (
	"nebula-go/pkg/gbc/memory/segments"
)

// FIXME: Implementatin directly copied from C++ version without checking, need
//  to verify this is the expected behavior.

type mbc5 struct {
	defaultImpl
}

func NewMBC5(rom, eram segments.Segment) MBC {
	return newMBCWrapper(rom, eram, &mbc5{
		defaultImpl: defaultImpl{
			rom:  rom,
			eram: eram,
		},
	})
}

func (m *mbc5) bankSelectorZone1(addr uint16, value uint8) error {
	var bank uint

	if addr < 0x3000 {
		bank = (m.rom.Bank() & 0x100) | uint(value)
	} else {
		bank = uint(value&0x1) << 8
		bank |= m.rom.Bank() & 0xFF
	}
	return m.rom.SelectBank(bank)
}

func (m *mbc5) bankSelectorZone2(value uint8) error {
	return m.eram.SelectBank(uint(value & 0x0F))
}

func (m *mbc5) bankModeSelect(value uint8) error {
	return nil
}
