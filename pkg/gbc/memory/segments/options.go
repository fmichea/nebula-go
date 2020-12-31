package segments

// Option configure segments with various behaviors.
type Option func(*segment) error

// WithBanks configures a given segment with a certain number of banks.
func WithBanks(count uint) Option {
	return func(s *segment) error {
		return s.banksCfg.setBankCount(count)
	}
}

// WithPunnedBank0 splits the segments in two equal parts, where the first half
// is for bank 0 exclusively, and the second half is for bank 1-N. More than
// one bank, configured using WithBanks, is mandatory with this option.
func WithPinnedBank0() Option {
	return func(s *segment) error {
		s.banksCfg.makeBank0Pinned()
		return nil
	}
}

// WithInitialData adds the given initial data to the segment. This may be used
// to initialize data from ROM dump.
func WithInitialData(buffer []uint8) Option {
	return func(s *segment) error {
		s.initialBuffer = buffer
		return nil
	}
}

// WithMirrorMapping maps startAddr:endAddr to this segment's address space.
// Both address spaces must be the same size.
func WithMirrorMapping(startAddr, endAddr uint16) Option {
	return func(s *segment) error {
		s.mirrorRanges = append(s.mirrorRanges, AddressRange{
			Start: startAddr,
			End:   endAddr,
		})
		return nil
	}
}
