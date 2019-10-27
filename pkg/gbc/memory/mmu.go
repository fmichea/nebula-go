package memory

import (
	"io"

	"nebula-go/pkg/gbc/memory/cartridge"
	"nebula-go/pkg/gbc/memory/lib"
	"nebula-go/pkg/gbc/memory/mbcs"
	"nebula-go/pkg/gbc/memory/segments"
)

type Registers struct{}

type MMU interface {
	Registers() *Registers

	ReadByte(addr uint16) (uint8, error)
	WriteByte(addr uint16, value uint8) error

	ReadDByte(addr uint16) (uint16, error)
	WriteDByte(addr, value uint16) error
}

type mmu struct {
	s   map[string]segments.Segment
	mbc mbcs.MBC

	regs *Registers
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
	results := map[string]segments.Segment{}

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
		results[cfg.name] = s
	}

	mbc := cr.MBCSelector.GetMBC(results["ROM"], results["ERAM"])
	if mbc == nil {
		return nil, lib.ErrMBCNotImplemented
	}

	result := &mmu{
		s:   results,
		mbc: mbc,
	}

	return result, nil
}

func (m *mmu) Registers() *Registers {
	return m.regs
}

func (m *mmu) ReadByte(addr uint16) (uint8, error) {
	return m.readByteInternal(addr)
}

func (m *mmu) ReadDByte(addr uint16) (uint16, error) {
	var result uint16

	if value, err := m.readByteInternal(addr + 1); err == nil {
		result |= uint16(value) << 8
	} else {
		return 0, err
	}

	if value, err := m.readByteInternal(addr); err == nil {
		result |= uint16(value)
	} else {
		return 0, err
	}

	return result, nil
}

func (m *mmu) WriteByte(addr uint16, value uint8) error {
	return nil // FIXME
}

func (m *mmu) WriteDByte(addr, value uint16) error {
	return nil // FIXME
}

func (m *mmu) readByteInternal(addr uint16) (uint8, error) {
	ptr, err := m.realBytePtr(lib.AccessTypeRead, addr, 0)
	if err != nil {
		return 0, err
	}

	if ptr != nil {
		result := m.readByteMasking(addr, *ptr)
		// FIXME: add memory watch notification here.
		return result, nil
	}

	return 0, lib.ErrInvalidRead
}

func (m *mmu) readByteMasking(addr uint16, value uint8) uint8 {
	// FIXME: handle NR10...NRXX masks here.
	return value
}

func (m *mmu) realBytePtr(accessType lib.AccessType, addr uint16, value uint8) (*uint8, error) {
	// FIXME: Handle BGPD here.

	// FIXME: return error when read on MBC with uninitialized? Cannot really happen?
	if m.mbc.ContainsAddress(addr) {
		return m.mbc.BytePtr(accessType, addr, value)
	}

	if segment := m.getSegment(addr); segment != nil {
		return segment.BytePtr(addr), nil
	}

	return nil, nil
}

func (m *mmu) getSegment(addr uint16) segments.Segment {
	// ROM and ERAM excluded because they are accessed through MBC.
	for _, segment := range m.s {
		if segment.ContainsAddress(addr) {
			return segment
		}
	}
	return nil
}
