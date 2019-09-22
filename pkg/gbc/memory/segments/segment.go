package segments

type Segment interface {
	ContainsAddress(addr uint16) bool
	BytePtr(addr uint16) *uint8

	SelectBank(bank uint) error
}

type Option func(*segment) error

func WithBanks(count uint) Option {
	return func(s *segment) error {
		return s.banksCfg.setBankCount(count)
	}
}

func WithPinnedBank0() Option {
	return func(s *segment) error {
		s.banksCfg.makeBank0Pinned()
		return nil
	}
}

func WithInitialData(buffer []uint8) Option {
	return func(s *segment) error {
		s.initialBuffer = buffer
		return nil
	}
}

func WithMirrorMapping(startAddr, endAddr uint16) Option {
	return func(s *segment) error {
		s.mirrorRanges = append(s.mirrorRanges, addressRange{
			start: startAddr,
			end:   endAddr,
		})
		return nil
	}
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
		addressRange: addressRange{
			start: startAddr,
			end:   endAddr,
		},

		banksCfg: newBanksConfig(startAddr, endAddr),
	}

	for _, opt := range opts {
		if err := opt(result); err != nil {
			return nil, err
		}
	}

	if err := result.initializeAndValidate(); err != nil {
		return nil, err
	}

	return result, nil
}

type segment struct {
	addressRange addressRange

	banksCfg banksConfig

	buffer        []uint8
	initialBuffer []uint8

	mirrorRanges []addressRange
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

func (s *segment) BytePtr(addr uint16) *uint8 {
	if s.banksCfg.containsAddress(addr) {
		return s.bytePtrForBankAndOffset(s.banksCfg.asOffset(addr), s.banksCfg.current)
	} else if s.addressRange.containsAddress(addr) {
		return s.bytePtrForBankAndOffset(s.addressRange.asOffset(addr), 0)
	}

	for _, mr := range s.mirrorRanges {
		if mr.containsAddress(addr) {
			return s.BytePtr(mr.transposeAddress(s.addressRange, addr))
		}
	}

	return nil
}

func (s *segment) SelectBank(bank uint) error {
	return s.banksCfg.selectBank(bank)
}

// Due to the options model, some of the initialization must be done after all options functions
// have been executed. This is what this function does.
func (s *segment) initializeAndValidate() error {
	// Initialize and validate the banks configuration.
	if err := s.banksCfg.initializeAndValidate(); err != nil {
		return err
	}

	// Create the internal buffer of the proper size.
	bufferSize := s.banksCfg.count * s.banksCfg.sizePerBank()
	s.buffer = make([]uint8, bufferSize)

	// Copying the initial buffer, if we did not copy the full buffer then the internal buffer is smaller, therefore
	// the buffers were incompatible.
	if n := copy(s.buffer, s.initialBuffer); n != len(s.initialBuffer) {
		return ErrBufferIncompatible
	}

	// Mirrored ranges must fit into the segment, the mirrored range may cross banked boundaries.
	for _, mr := range s.mirrorRanges {
		if s.addressRange.size() < mr.size() {
			return ErrInvalidMirrorRange
		}
	}

	return nil
}

func (s *segment) bytePtrForBankAndOffset(addr uint16, bank uint) *uint8 {
	offset := uint(addr) + s.banksCfg.sizePerBank()*bank
	return &s.buffer[offset]
}
