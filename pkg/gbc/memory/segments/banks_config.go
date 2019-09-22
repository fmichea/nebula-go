package segments

type banksConfig struct {
	segmentAddressRange addressRange
	bankedAddressRange  addressRange

	current uint
	count   uint

	bank0Pinned bool
}

func newBanksConfig(startAddr, endAddr uint16) banksConfig {
	return banksConfig{
		segmentAddressRange: addressRange{
			start: startAddr,
			end:   endAddr,
		},
		bankedAddressRange: addressRange{
			start: startAddr,
			end:   endAddr,
		},

		current: 0,
		count:   1,

		bank0Pinned: false,
	}
}

func (c *banksConfig) setBankCount(count uint) error {
	if count == 0 {
		return ErrBankCountInvalid
	}
	c.count = count
	return nil
}

func (c *banksConfig) makeBank0Pinned() {
	// Mark banks config with bank0 pinned.
	c.bank0Pinned = true

	// Move banks to the second half of the segment address range.
	c.bankedAddressRange.start = c.segmentAddressRange.start
	c.bankedAddressRange.start += c.segmentAddressRange.size() / 2
}

func (c *banksConfig) asOffset(addr uint16) uint16 {
	return c.bankedAddressRange.asOffset(addr)
}

func (c *banksConfig) containsAddress(addr uint16) bool {
	return c.bankedAddressRange.containsAddress(addr)
}

func (c *banksConfig) sizePerBank() uint {
	return uint(c.bankedAddressRange.size())
}

func (c *banksConfig) initializeAndValidate() error {
	if c.bank0Pinned && c.count == 1 {
		return ErrCannotPin0WithOneBank
	}

	if c.bank0Pinned {
		c.current = 1
	} else {
		c.current = 0
	}
	return nil
}

func (c *banksConfig) selectBank(bank uint) error {
	if c.count <= bank || (bank == 0 && c.bank0Pinned) {
		return ErrBankUnavailable
	}

	c.current = bank
	return nil
}
