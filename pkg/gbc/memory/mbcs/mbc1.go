package mbcs

import (
	"go.uber.org/multierr"

	"nebula-go/pkg/gbc/memory/segments"
)

type mbc1mode int

const (
	romBankingMode mbc1mode = iota
	ramBankingMode
)

type mbc1 struct {
	defaultImpl

	mode             mbc1mode
	secondaryBanking uint8
}

func newMBC1(rom, eram segments.Segment) MBC {
	return newMBCWrapper(rom, eram, &mbc1{
		defaultImpl: defaultImpl{
			rom:  rom,
			eram: eram,
		},

		mode: romBankingMode,
	})
}

func (i *mbc1) bankSelectorZone1(addr uint16, value uint8) error {
	var romBank uint8

	if i.mode == romBankingMode {
		romBank = uint8(i.rom.Bank()) & 0xE0
	}

	value &= 0x1F
	if value == 0 {
		value += 1
	}

	return i.rom.SelectBank(uint(romBank | value))
}

func (i *mbc1) bankSelectorZone2(value uint8) (err error) {
	i.secondaryBanking = value & 3

	switch i.mode {
	case romBankingMode:
		err = i.selectUpperROMBits(i.secondaryBanking)

	case ramBankingMode:
		err = i.selectRAMBank(i.secondaryBanking)
	}
	return
}

func (i *mbc1) bankModeSelect(value uint8) error {
	if (value & 0x01) == 1 {
		i.mode = ramBankingMode

		return multierr.Combine(
			i.selectUpperROMBits(0),
			i.selectRAMBank(i.secondaryBanking),
		)
	}

	i.mode = romBankingMode

	return multierr.Combine(
		i.selectUpperROMBits(i.secondaryBanking),
		i.selectRAMBank(0),
	)
}

func (i *mbc1) selectUpperROMBits(value uint8) error {
	romBank := uint8(i.rom.Bank()) & 0x1F
	romBank |= value << 5
	return i.rom.SelectBank(uint(romBank))
}

func (i *mbc1) selectRAMBank(value uint8) error {
	return i.eram.SelectBank(uint(value))
}
