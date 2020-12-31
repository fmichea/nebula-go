package mbcs

import (
	"nebula-go/pkg/gbc/memory/lib"
	"nebula-go/pkg/gbc/memory/segments"
)

type mbcImpl interface {
	readRAMAddress(addr uint16) (uint8, error)
	readRAMAddressSlice(addr uint16, count uint) ([]uint8, error)

	writeRAMAddress(addr uint16, value uint8) error

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

func (w *mbcWrapper) AddressRanges() []segments.AddressRange {
	var result []segments.AddressRange

	result = append(result, w.rom.AddressRanges()...)
	result = append(result, w.eram.AddressRanges()...)
	return result
}

func (w *mbcWrapper) ReadByte(addr uint16) (uint8, error) {
	if w.rom.ContainsAddress(addr) {
		return w.rom.ReadByte(addr)
	} else if w.eram.ContainsAddress(addr) {
		return w.impl.readRAMAddress(addr)
	}
	return 0, lib.ErrInvalidRead
}

func (w *mbcWrapper) ReadByteSlice(addr uint16, count uint) ([]uint8, error) {
	if w.rom.ContainsAddress(addr) {
		return w.rom.ReadByteSlice(addr, count)
	} else if w.eram.ContainsAddress(addr) {
		return w.impl.readRAMAddressSlice(addr, count)
	}
	return nil, lib.ErrInvalidRead
}

func (w *mbcWrapper) WriteByte(addr uint16, value uint8) error {
	// External RAM is writeable.
	if w.eram.ContainsAddress(addr) {
		return w.impl.writeRAMAddress(addr, value)
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
	switch {
	case addr <= 0x1FFF:
		return w.impl.ramEnable(value)

	case 0x2000 <= addr && addr <= 0x3FFF:
		return w.impl.bankSelectorZone1(addr, value)

	case 0x4000 <= addr && addr <= 0x5FFF:
		return w.impl.bankSelectorZone2(value)

	case 0x6000 <= addr && addr <= 0x7FFF:
		return w.impl.bankModeSelect(value)

	default:
		return lib.ErrInvalidWrite
	}
}

func (w *mbcWrapper) WriteByteSlice(addr uint16, values []uint8) error {
	// TODO: there is nothing in specification disallowing this, not currently used but could be implemented.
	return ErrMBCSliceOperatioInvalid
}

func (w *mbcWrapper) ByteHook(addr uint16) (*uint8, error) {
	return nil, ErrMBCHookInvalid
}

type defaultImpl struct {
	rom  segments.Segment
	eram segments.Segment

	ramEnabled bool
}

func (i *defaultImpl) readRAMAddress(addr uint16) (uint8, error) {
	if !i.ramEnabled {
		return 0, ErrRAMUnavailable
	}
	return i.eram.ReadByte(addr)
}

func (i *defaultImpl) readRAMAddressSlice(addr uint16, count uint) ([]uint8, error) {
	if !i.ramEnabled {
		return nil, ErrRAMUnavailable
	}
	return i.eram.ReadByteSlice(addr, count)
}

func (i *defaultImpl) writeRAMAddress(addr uint16, value uint8) error {
	if !i.ramEnabled {
		return ErrRAMUnavailable
	}
	return i.eram.WriteByte(addr, value)
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
