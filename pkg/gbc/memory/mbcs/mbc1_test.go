package mbcs

func (s *unitTestSuite) TestMBC1_ContainsAddr() {
	s.testContainsAddress(s.mbc1)
}

func (s *unitTestSuite) TestMBC1_ROMRead() {
	value, err := s.mbc1.ReadByte(0x0000)
	s.Require().NoError(err)
	s.Equal(uint8(0), value)

	value, err = s.mbc1.ReadByte(0x4000)
	s.Require().NoError(err)
	s.Equal(uint8(1), value)
}

func (s *unitTestSuite) TestMBC1_ERAMRead() {
	// Before RAM is enabled, read is not available.
	_, err := s.mbc1.ReadByte(0xA000)
	s.Require().Equal(ErrRAMUnavailable, err)

	// RAM enable.
	err = s.mbc1.WriteByte(0x1FFF, 0x0A)
	s.Require().NoError(err)

	// Bank 0 read in ERAM.
	value, err := s.mbc1.ReadByte(0xA000)
	s.Require().NoError(err)
	s.Equal(uint8(0), value)

	// RAM disable.
	err = s.mbc1.WriteByte(0x1FFF, 0x00)
	s.Require().NoError(err)

	// RAM read is not available anymore.
	_, err = s.mbc1.ReadByte(0xA000)
	s.Require().Equal(ErrRAMUnavailable, err)
}

func (s *unitTestSuite) TestMBC1_ERAMWrite() {
	// Before RAM is enabled, read is not available.
	err := s.mbc1.WriteByte(0xA000, 0)
	s.Require().Equal(ErrRAMUnavailable, err)

	// RAM enable.
	err = s.mbc1.WriteByte(0x1FFF, 0x0A)
	s.NoError(err)

	// Bank 0 read in ERAM.
	err = s.mbc1.WriteByte(0xA000, 0xF)
	s.NoError(err)

	value, err := s.mbc1.ReadByte(0xA000)
	s.Require().NoError(err)
	s.Equal(uint8(0x0F), value)

	// RAM disable.
	err = s.mbc1.WriteByte(0x1FFF, 0x00)
	s.Require().NoError(err)

	// No access to RAM anymore.
	err = s.mbc1.WriteByte(0xA000, 0)
	s.Require().Equal(ErrRAMUnavailable, err)
}

func (s *unitTestSuite) TestMBC1_ROMBankSelect() {
	selectLow := func(value uint8) {
		// Write value in the rom bank select zone.
		err := s.mbc1.WriteByte(0x2000, value)
		s.Require().NoError(err, "failed to select bank %#x (low bits)", value)
	}

	selectHigh := func(value uint8) {
		err := s.mbc1.WriteByte(0x4000, value)
		s.Require().NoError(err, "failed to select bank %#x (high bits)", value)
	}

	checkBankSelected := func(expectedBank uint8) {
		// Now check the expected bank was selected.
		value, err := s.mbc1.ReadByte(0x4000)
		s.Require().NoError(err)
		s.Equal(expectedBank, value)
	}

	// Trying to select bank 0 translates selection to bank 1.
	selectLow(0x00)
	checkBankSelected(0x01)

	// Selecting bank 2, 3 and 0xF works as expected.
	for _, value := range []uint8{0x02, 0x03, 0x0F} {
		selectLow(value)
		checkBankSelected(value)
	}

	// Selecting more than 5 bits here only selects the lower 5.
	selectLow(0xFF)
	checkBankSelected(0x1F)

	// Selecting upper bits is allowed.
	selectHigh(0x01)
	checkBankSelected(0x3F)

	// Selecting bank 0 low still has same 0->1 behavior.
	selectLow(0)
	checkBankSelected(0x21)

	// Only 3 bits are used.
	selectHigh(0xFF)
	checkBankSelected(0x61)
}

func (s *unitTestSuite) TestMBC1_RAMBankSelect() {
	// Enable RAM and set banking mode to RAM banking.
	err := s.mbc1.WriteByte(0x0000, 0x0A)
	s.Require().NoError(err)

	err = s.mbc1.WriteByte(0x6000, 0x01)
	s.Require().NoError(err)

	selectBank := func(value uint8) {
		err := s.mbc1.WriteByte(0x4000, value)
		s.Require().NoError(err, "failed to select bank %#x", value)
	}

	checkBank := func(expectedBank uint8) {
		// Now check the expected bank was selected.
		value, err := s.mbc1.ReadByte(0xB000)
		s.Require().NoError(err)
		s.Equal(expectedBank, value)
	}

	values := map[uint8]uint8{
		0x00: 0x00,
		0x01: 0x01,
		0x02: 0x02,
		0x03: 0x03,
		0xFF: 0x03,
	}
	for bankWanted, bankSelected := range values {
		selectBank(bankWanted)
		checkBank(bankSelected)
	}
}

func (s *unitTestSuite) TestMBC1_SwitchingBetweenRAMAndROMBanking() {
	// Enable RAM and set banking mode to RAM banking.
	err := s.mbc1.WriteByte(0x0000, 0x0A)
	s.Require().NoError(err)

	// Select the highest bank we can in ROM: 0x7F
	err = s.mbc1.WriteByte(0x2000, 0xFF)
	s.Require().NoError(err)

	err = s.mbc1.WriteByte(0x4000, 0xFF)
	s.Require().NoError(err)

	// Ensure the bank got selected correctly.
	value, err := s.mbc1.ReadByte(0x4000)
	s.Require().NoError(err)
	s.Equal(uint8(0x7F), value)

	value, err = s.mbc1.ReadByte(0xA000)
	s.Require().NoError(err)
	s.Equal(uint8(0x00), value)

	// Now switch to RAM banking, which should reduce to bank 0x1F at most.
	err = s.mbc1.WriteByte(0x6000, 0x01)
	s.Require().NoError(err)

	// Selected ROM bank is 0x1F!
	value, err = s.mbc1.ReadByte(0x4000)
	s.Require().NoError(err)
	s.Equal(uint8(0x1F), value)

	// RAM selected bank is now 3.
	value, err = s.mbc1.ReadByte(0xA000)
	s.Require().NoError(err)
	s.Equal(uint8(0x03), value)

	// Going back to ROM banking switches back the selector.
	err = s.mbc1.WriteByte(0x6000, 0x00)
	s.Require().NoError(err)

	// Selected ROM bank is back to 0x7F.
	value, err = s.mbc1.ReadByte(0x4000)
	s.NoError(err)
	s.Equal(uint8(0x7F), value)

	// RAM selected bank is now 0.
	value, err = s.mbc1.ReadByte(0xA000)
	s.Require().NoError(err)
	s.Equal(uint8(0x00), value)
}
