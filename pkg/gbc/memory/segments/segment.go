package segments

import (
	"nebula-go/pkg/gbc/memory/lib"
)

type SegmentBase interface {
	AddressRanges() []AddressRange
	ContainsAddress(addr uint16) bool

	ReadByte(addr uint16) (uint8, error)
	ReadByteSlice(addr uint16, count uint) ([]uint8, error)

	WriteByte(addr uint16, value uint8) error
	WriteByteSlice(addr uint16, values []uint8) error

	ByteHook(addr uint16) (*uint8, error)
}

type Segment interface {
	SegmentBase

	Bank() uint
	BankCount() uint
	SelectBank(bank uint) error
}

// New creates a new segment. SegmentOptions can be used to change the behavior of the segment.
//
// By default, there is no banking and the segment looks like this:
// +-------------------------+
// |          Bank 0         |
// +-------------------------+
//
// When using banking (WithBanks), the segment will be in the following structure (N = bankCount - 1):
// +-------------------------+
// |         Bank 0-N        |
// +-------------------------+
//
// If the first bank is pinned (WithPinnedBank0), the segment is split in two halves. The bank indexed at 0 is
// available at the first half of the segment, and the other banks are available over the second half of the
// segment. It would have the following structure (N = bankCount - 1):
// +------------+------------+
// |   Bank 0   |  Bank 1-N  |
// +------------+------------+
func New(startAddr, endAddr uint16, opts ...Option) (*segment, error) {
	result := &segment{
		addressRange: AddressRange{
			Start: startAddr,
			End:   endAddr,
		},

		banksCfg: newBanksConfig(startAddr, endAddr),
	}

	for _, opt := range opts {
		if err := opt(result); err != nil {
			return nil, err
		}
	}

	if err := result.validateAndInitialize(); err != nil {
		return nil, err
	}

	return result, nil
}

type segment struct {
	addressRange AddressRange

	banksCfg banksConfig

	buffers       [][]uint8
	initialBuffer []uint8

	mirrorRanges []AddressRange
}

func (s *segment) AddressRanges() []AddressRange {
	var result []AddressRange

	result = append(result, s.addressRange)
	for _, mr := range s.mirrorRanges {
		result = append(result, mr)
	}
	return result
}

func (s *segment) ContainsAddress(addr uint16) bool {
	if s.addressRange.containsAddress(addr) {
		return true
	}

	for _, mr := range s.mirrorRanges {
		if mr.containsAddress(addr) {
			return true
		}
	}

	return false
}

func (s *segment) ReadByte(addr uint16) (uint8, error) {
	ptr, err := s.bytePtr(addr)
	if err != nil {
		return 0, err
	}
	return *ptr, nil
}

func (s *segment) ReadByteSlice(addr uint16, count uint) ([]uint8, error) {
	bank, offset, err := s.byteOffset(addr)
	if err != nil {
		return nil, err
	}

	if !s.addressRange.hasCapacityFromAddress(addr, count) {
		return nil, ErrSegmentTooSmall
	}

	bankData := s.buffers[bank]
	bankDataSizeLeftFromOffset := uint(len(bankData)) - offset

	if count <= bankDataSizeLeftFromOffset {
		data := bankData[offset : offset+count]
		return data, nil
	}

	data := append(
		bankData[offset:],
		s.buffers[s.banksCfg.current][:count-bankDataSizeLeftFromOffset]...,
	)
	return data, nil
}

func (s *segment) WriteByte(addr uint16, value uint8) error {
	ptr, err := s.bytePtr(addr)
	if err != nil {
		return err
	}

	*ptr = value
	return nil
}

func (s *segment) WriteByteSlice(addr uint16, values []uint8) error {
	count := uint(len(values))

	bank, offset, err := s.byteOffset(addr)
	if err != nil {
		return err
	}

	if !s.addressRange.hasCapacityFromAddress(addr, count) {
		return ErrSegmentTooSmall
	}

	bankData := s.buffers[bank]
	bankDataSizeLeftFromOffset := uint(len(bankData)) - offset

	if count <= bankDataSizeLeftFromOffset {
		for idx, value := range values {
			bankData[offset+uint(idx)] = value
		}
	} else {
		for idx := uint(0); idx < bankDataSizeLeftFromOffset; idx++ {
			bankData[offset+idx] = values[idx]
		}

		for idx, value := range values[bankDataSizeLeftFromOffset:] {
			s.buffers[s.banksCfg.current][idx] = value
		}
	}
	return nil
}

func (s *segment) ByteHook(addr uint16) (*uint8, error) {
	if s.banksCfg.containsAddress(addr) && s.banksCfg.isBanked(addr) {
		return nil, ErrInvalidHookInBank
	}
	return s.bytePtr(addr)
}

func (s *segment) byteOffset(addr uint16) (uint, uint, error) {
	if s.banksCfg.containsAddress(addr) {
		return s.banksCfg.current, s.banksCfg.asOffset(addr), nil
	} else if s.addressRange.containsAddress(addr) {
		return 0, s.addressRange.asOffset(addr), nil
	}

	for _, mr := range s.mirrorRanges {
		return s.byteOffset(mr.transposeAddress(s.addressRange, addr))
	}

	return 0, 0, lib.ErrInvalidSegmentAddr
}

func (s *segment) bytePtr(addr uint16) (*uint8, error) {
	bank, offset, err := s.byteOffset(addr)
	if err != nil {
		return nil, err
	}
	return &s.buffers[bank][offset], nil
}

func (s *segment) Bank() uint {
	return s.banksCfg.current
}

func (s *segment) BankCount() uint {
	return s.banksCfg.count
}

func (s *segment) SelectBank(bank uint) error {
	return s.banksCfg.selectBank(bank)
}

// Due to the options model, some of the initialization must be done after all options functions
// have been executed. This is what this function does.
func (s *segment) validateAndInitialize() error {
	// Initialize and validate the banks configuration.
	if err := s.banksCfg.validateAndInitialize(); err != nil {
		return err
	}

	sizePerBank := s.banksCfg.sizePerBank()

	// Create the internal buffer for each bank within this segment.
	buffer := make([]uint8, s.banksCfg.count*sizePerBank)

	// Copying the initial buffer, if we did not copy the full buffer then the internal buffer is smaller, therefore
	// the buffers were incompatible.
	if n := copy(buffer, s.initialBuffer); n != len(s.initialBuffer) {
		return ErrBufferIncompatible
	}

	// Now we split the buffer into each bank.
	s.buffers = make([][]uint8, s.banksCfg.count)
	for bankIdx := uint(0); bankIdx < s.banksCfg.count; bankIdx++ {
		s.buffers[bankIdx] = buffer[bankIdx*sizePerBank : (bankIdx+1)*sizePerBank]
	}

	// Mirrored ranges must fit into the segment, the mirrored range may cross banked boundaries.
	for _, mr := range s.mirrorRanges {
		if s.addressRange.size() < mr.size() {
			return ErrInvalidMirrorRange
		}
	}

	return nil
}
