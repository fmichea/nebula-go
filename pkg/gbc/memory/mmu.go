package memory

import (
	"io"

	"nebula-go/pkg/gbc/memory/cartridge"
	"nebula-go/pkg/gbc/memory/lib"
	"nebula-go/pkg/gbc/memory/mbcs"
	"nebula-go/pkg/gbc/memory/segments"
)

type MMU interface {
	lib.MemoryIO

	Cartridge() *cartridge.ROM
	Registers() *Registers
}

type mmu struct {
	segmentsMapping     map[string]segments.Segment
	segmentsAddrMapping []segments.SegmentBase
	mbc                 mbcs.MBC

	hooks map[uint16]lib.Hook

	cartridge *cartridge.ROM
	regs      *Registers
}

func NewMMUFromFile(out io.Writer, filename string) (MMU, error) {
	cr, err := cartridge.Load(out, filename)
	if err != nil {
		return nil, err
	}
	return NewMMUFromCartridge(cr)
}

func NewMMUFromCartridge(cr *cartridge.ROM) (MMU, error) {
	return newMMUFromCartridge(cr)
}

func newMMUFromCartridge(cr *cartridge.ROM) (*mmu, error) {
	var err error

	segmentsMapping := map[string]segments.Segment{}

	configs := []struct {
		name      string
		startAddr uint16
		endAddr   uint16
		options   []segments.Option
	}{
		{
			name:      "ROM",
			startAddr: 0x0000,
			endAddr:   0x7FFF,
			options: []segments.Option{
				segments.WithBanks(cr.Size.BankCount()),
				segments.WithPinnedBank0(),
				segments.WithInitialData(cr.Data),
			},
		},
		{
			name:      "VRAM",
			startAddr: 0x8000,
			endAddr:   0x9FFF,
			options: []segments.Option{
				segments.WithBanks(cr.Type.VRAMBankCount()),
			},
		},
		{
			name:      "ERAM",
			startAddr: 0xA000,
			endAddr:   0xBFFF,
			options: []segments.Option{
				segments.WithBanks(cr.RAMSize.BankCount()),
			},
		},
		{
			name:      "WRAM",
			startAddr: 0xC000,
			endAddr:   0xDFFF,
			options: []segments.Option{
				segments.WithBanks(cr.Type.WRAMBankCount()),
				segments.WithPinnedBank0(),
				segments.WithMirrorMapping(0xE000, 0xFDFF), // ECHO segment.
			},
		},
		{
			name:      "OAM",
			startAddr: 0xFE00,
			endAddr:   0xFE9F,
		},
		{
			name:      "IO_PORTS",
			startAddr: 0xFF00,
			endAddr:   0xFF7F,
		},
		{
			name:      "HRAM",
			startAddr: 0xFF80,
			endAddr:   0xFFFF,
		},
	}

	for _, cfg := range configs {
		s, err := segments.New(cfg.startAddr, cfg.endAddr, cfg.options...)
		if err != nil {
			return nil, err
		}
		segmentsMapping[cfg.name] = s
	}

	mbc := cr.MBCSelector.GetMBC(segmentsMapping["ROM"], segmentsMapping["ERAM"])
	if mbc == nil {
		return nil, lib.ErrMBCNotImplemented
	}

	segmentsAddrMapping := make([]segments.SegmentBase, 0x10000)
	for _, s := range []segments.SegmentBase{
		mbc,
		segmentsMapping["VRAM"],
		segmentsMapping["WRAM"],
		segmentsMapping["OAM"],
		segmentsMapping["IO_PORTS"],
		segmentsMapping["HRAM"],
	} {
		for _, addressRange := range s.AddressRanges() {
			for addr := uint32(addressRange.Start); addr <= uint32(addressRange.End); addr++ {
				segmentsAddrMapping[addr] = s
			}
		}
	}

	result := &mmu{
		cartridge: cr,

		segmentsMapping:     segmentsMapping,
		segmentsAddrMapping: segmentsAddrMapping,
		mbc:                 mbc,

		hooks: map[uint16]lib.Hook{},
	}

	result.regs, err = InitializeRegisters(result, cr, segmentsMapping["VRAM"], segmentsMapping["WRAM"])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *mmu) Cartridge() *cartridge.ROM {
	return m.cartridge
}

func (m *mmu) Registers() *Registers {
	return m.regs
}

func (m *mmu) getSegmentForAddress(addr uint16) (segments.SegmentBase, error) {
	segment := m.segmentsAddrMapping[addr]
	if segment == nil {
		return nil, lib.ErrNoSegmentAtAddr
	}
	return segment, nil
}
