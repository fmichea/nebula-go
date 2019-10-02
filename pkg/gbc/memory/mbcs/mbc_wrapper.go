package mbcs

import (
	"errors"

	"nebula-go/pkg/gbc/memory/lib"
	"nebula-go/pkg/gbc/memory/segments"
)

type mbcImpl interface {
	readRAMAddress(addr uint16) (*uint8, error)
	writeRAMAddress(addr uint16, value uint8) (*uint8, error)

	ramEnable(value uint8) error
	bankSelectorZone1(addr uint16, value uint8) error
	bankSelectorZone2(value uint8) error
	bankModeSelect(value uint8) error
}

type mbcWrapper struct {
	rom  segments.Segment
	eram segments.Segment

	impl mbcImpl
}

func newMBCWrapper(rom, eram segments.Segment, impl mbcImpl) MBC {
	return &mbcWrapper{
		rom:  rom,
		eram: eram,

		impl: impl,
	}
}

func (w *mbcWrapper) ContainsAddress(addr uint16) bool {
	return w.rom.ContainsAddress(addr) || w.eram.ContainsAddress(addr)
}

var (
	ErrInvalidRead    = errors.New("invalid read received by MBC")
	ErrInvalidWrite   = errors.New("invalid write received by MBC")
	ErrRAMUnavailable = errors.New("RAM is not available at the moment")
)

func (w *mbcWrapper) BytePtr(accessType lib.AccessType, addr uint16, value uint8) (ptr *uint8, err error) {
	switch accessType {
	case lib.AccessTypeRead:
		if w.rom.ContainsAddress(addr) {
			ptr = w.rom.BytePtr(addr)
		} else if w.eram.ContainsAddress(addr) {
			ptr, err = w.impl.readRAMAddress(addr)
		} else {
			err = ErrInvalidRead
		}

	case lib.AccessTypeWrite:
		// External RAM is writeable, we return a pointer to it.
		if w.eram.ContainsAddress(addr) {
			ptr, err = w.impl.writeRAMAddress(addr, value)
			return
		}

		// Bank Selector:
		//   Writing to the ROM is used to controller the MBC behavior.
		//
		// Zones:
		//   0000h-1FFFh: RAM Enable, used to protect the RAM during
		//                power-down, not implemented yet.
		//   2000h-3FFFh: ROM Bank Selector, different for every controller.
		//   4000h-5FFFh: ROM/RAM Bank Selector, different for every controller.
		//   6000h-7FFFh: MBC mode controller, for MBC1 and MBC3.
		if addr <= 0x1FFF {
			err = w.impl.ramEnable(value)
		} else if 0x2000 <= addr && addr <= 0x3FFF {
			err = w.impl.bankSelectorZone1(addr, value)
		} else if 0x4000 <= addr && addr <= 0x5FFF {
			err = w.impl.bankSelectorZone2(value)
		} else if 0x6000 <= addr && addr <= 0x7FFF {
			err = w.impl.bankModeSelect(value)
		} else {
			err = ErrInvalidWrite
		}
	}
	return
}

type defaultImpl struct {
	rom  segments.Segment
	eram segments.Segment

	ramEnabled bool
}

func (i *defaultImpl) readRAMAddress(addr uint16) (*uint8, error) {
	if !i.ramEnabled {
		return nil, ErrRAMUnavailable
	}
	return i.eram.BytePtr(addr), nil
}

func (i *defaultImpl) writeRAMAddress(addr uint16, value uint8) (*uint8, error) {
	if !i.ramEnabled {
		return nil, ErrRAMUnavailable
	}
	return i.eram.BytePtr(addr), nil
}

func (i *defaultImpl) ramEnable(value uint8) error {
	i.ramEnabled = (value & 0x0F) == 0x0A
	return nil
}

//func (i *defaultImpl) bankSelectorZone1(addr uint16, value uint8) error {
//	return i.rom.SelectBank(uint(value))
//}
//
//func (i *defaultImpl) bankSelectorZone2(value uint8) error {
//	return i.eram.SelectBank(uint(value))
//}
//
//func (i *defaultImpl) bankModeSelect(value uint8) error {
//	return nil
//}
