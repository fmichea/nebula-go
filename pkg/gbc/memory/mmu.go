package memory

import (
	"io"

	"nebula-go/pkg/gbc/memory/cartridge"
	"nebula-go/pkg/gbc/memory/lib"
	"nebula-go/pkg/gbc/memory/mbcs"
	"nebula-go/pkg/gbc/memory/segments"
)

type MMU struct {
	s   map[string]segments.Segment
	mbc mbcs.MBC
}

func NewMMU(out io.Writer, filename string) (*MMU, error) {
	cr, err := cartridge.Load(out, filename)
	if err != nil {
		return nil, err
	}

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

	result := &MMU{
		s:   results,
		mbc: mbc,
	}

	return result, nil
}

func (m *MMU) ReadByte(addr uint16) uint8 {
	return m.readByteInternal(addr)
}

func (m *MMU) ReadDByte(addr uint16) uint16 {
	result := uint16(m.readByteInternal(addr+1)) << 8
	result |= uint16(m.readByteInternal(addr))
	return result
}

func (m *MMU) readByteInternal(addr uint16) uint8 {
	return 0
}

func (m *MMU) realBytePtr(addr uint16) *uint8 {
	// FIXME: Handle BGPD here.

	// FIXME: return error when read on MBC with uninitialized? Cannot really happen?
	if m.mbc.ContainsAddress(addr) {
		ptr, _ := m.mbc.BytePtr(lib.AccessTypeRead, addr, 0) // FIXME
		return ptr
	}

	if segment := m.getSegment(addr); segment != nil {
		return segment.BytePtr(addr)
	}
	return nil
}

func (m *MMU) getSegment(addr uint16) segments.Segment {
	// ROM and ERAM excluded because they are accessed through MBC.
	for _, segment := range m.s {
		if segment.ContainsAddress(addr) {
			return segment
		}
	}
	return nil
}
