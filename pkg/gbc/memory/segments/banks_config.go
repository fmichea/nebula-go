package segments

type banksConfig struct {
	segmentAddressRange AddressRange
	bankedAddressRange  AddressRange

	current uint
	count   uint

	bank0Pinned bool
}

func newBanksConfig(startAddr, endAddr uint16) banksConfig {
	return banksConfig{
		segmentAddressRange: AddressRange{
			Start: startAddr,
			End:   endAddr,
		},
		bankedAddressRange: AddressRange{
			Start: startAddr,
			End:   endAddr,
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
	c.bankedAddressRange.Start = c.segmentAddressRange.Start
	c.bankedAddressRange.Start += uint16(c.segmentAddressRange.size()) / 2
}

func (c *banksConfig) asOffset(addr uint16) uint {
	return c.bankedAddressRange.asOffset(addr)
}

func (c *banksConfig) containsAddress(addr uint16) bool {
	return c.bankedAddressRange.containsAddress(addr)
}

func (c *banksConfig) isBanked(addr uint16) bool {
	if c.bank0Pinned {
		return c.count != 2 && c.bankedAddressRange.containsAddress(addr)
	}
	return c.count != 1
}

func (c *banksConfig) sizePerBank() uint {
	return c.bankedAddressRange.size()
}

func (c *banksConfig) validateAndInitialize() error {
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
	if bank == 0 && c.bank0Pinned {
		bank = 1
	}

	if c.count <= bank /* || (bank == 0 && c.bank0Pinned) */ {
		return ErrBankUnavailable
	}

	c.current = bank
	return nil
}
