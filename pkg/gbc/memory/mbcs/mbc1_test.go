package mbcs

import (
	"nebula-go/pkg/gbc/memory/lib"
)

func (s *unitTestSuite) TestMBC1_ContainsAddr() {
	s.testContainsAddress(s.mbc1)
}

func (s *unitTestSuite) TestMBC1_ROMRead() {
	var err error
	var ptr *uint8

	ptr, err = s.mbc1.BytePtr(lib.AccessTypeRead, 0x0000, 0)
	s.NoError(err)
	s.NotNil(ptr)
	s.Equal(uint8(0), *ptr)

	ptr, err = s.mbc1.BytePtr(lib.AccessTypeRead, 0x4000, 0)
	s.NoError(err)
	s.NotNil(ptr)
	s.Equal(uint8(1), *ptr)
}

func (s *unitTestSuite) TestMBC1_ERAMRead() {
	var err error
	var ptr *uint8

	// Before RAM is enabled, read is not available.
	s.Nil(s.mbc1.BytePtr(lib.AccessTypeRead, 0xA000, 0x00))

	// RAM enable.
	_, err = s.mbc1.BytePtr(lib.AccessTypeWrite, 0x1FFF, 0x0A)
	s.NoError(err)

	// Bank 0 read in ERAM.
	ptr, err = s.mbc1.BytePtr(lib.AccessTypeRead, 0xA000, 0x00)
	s.NoError(err)
	s.NotNil(ptr)
	s.Equal(uint8(0), *ptr)

	// RAM disable.
	_, err = s.mbc1.BytePtr(lib.AccessTypeWrite, 0x1FFF, 0x00)
	s.NoError(err)

	// RAM read is not available anymore.
	ptr, err = s.mbc1.BytePtr(lib.AccessTypeRead, 0xA000, 0x00)
	s.Equal(ErrRAMUnavailable, err)
	s.Nil(ptr)
}

func (s *unitTestSuite) TestMBC1_ERAMWrite() {
	// Before RAM is enabled, read is not available.
	s.Nil(s.mbc1.BytePtr(lib.AccessTypeWrite, 0xA000, 0))

	// RAM enable.
	_, err := s.mbc1.BytePtr(lib.AccessTypeWrite, 0x1FFF, 0x0A)
	s.NoError(err)

	// Bank 0 read in ERAM.
	ptr, err := s.mbc1.BytePtr(lib.AccessTypeWrite, 0xA000, 0)
	s.NoError(err)
	s.NotNil(ptr)
	s.Equal(uint8(0), *ptr)
	*ptr = 0x0F

	ptr2, err := s.mbc1.BytePtr(lib.AccessTypeWrite, 0xA000, 0)
	s.NoError(err)
	s.Equal(ptr2, ptr)
	s.Equal(uint8(0x0F), *ptr2)

	// RAM disable.
	_, err = s.mbc1.BytePtr(lib.AccessTypeWrite, 0x1FFF, 0x00)
	s.NoError(err)

	// No access to RAM anymore.
	s.Nil(s.mbc1.BytePtr(lib.AccessTypeWrite, 0xA000, 0))
}

func (s *unitTestSuite) TestMBC1_ROMBankSelect() {
	selectLow := func(value uint8) {
		// Write value in the rom bank select zone.
		ptr, err := s.mbc1.BytePtr(lib.AccessTypeWrite, 0x2000, value)
		s.NoError(err, "failed to select bank %#x (low bits)", value)
		s.Nil(ptr)
	}

	selectHigh := func(value uint8) {
		ptr, err := s.mbc1.BytePtr(lib.AccessTypeWrite, 0x4000, value)
		s.NoError(err, "failed to select bank %#x (high bits)", value)
		s.Nil(ptr)
	}

	checkBankSelected := func(expectedBank uint8) {
		// Now check the expected bank was selected.
		ptr, err := s.mbc1.BytePtr(lib.AccessTypeRead, 0x4000, 0)
		s.NoError(err)
		s.NotNil(ptr)
		s.Equal(expectedBank, *ptr)
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
	_, err := s.mbc1.BytePtr(lib.AccessTypeWrite, 0x0000, 0x0A)
	s.NoError(err)

	_, err = s.mbc1.BytePtr(lib.AccessTypeWrite, 0x6000, 0x01)
	s.NoError(err)

	selectBank := func(value uint8) {
		ptr, err := s.mbc1.BytePtr(lib.AccessTypeWrite, 0x4000, value)
		s.NoError(err, "failed to select bank %#x", value)
		s.Nil(ptr)
	}

	checkBank := func(expectedBank uint8) {
		// Now check the expected bank was selected.
		ptr, err := s.mbc1.BytePtr(lib.AccessTypeRead, 0xB000, 0)
		s.NoError(err)
		s.NotNil(ptr)
		s.Equal(expectedBank, *ptr)
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
	_, err := s.mbc1.BytePtr(lib.AccessTypeWrite, 0x0000, 0x0A)
	s.NoError(err)

	// Select the highest bank we can in ROM: 0x7F
	_, err = s.mbc1.BytePtr(lib.AccessTypeWrite, 0x2000, 0xFF)
	s.NoError(err)

	_, err = s.mbc1.BytePtr(lib.AccessTypeWrite, 0x4000, 0xFF)
	s.NoError(err)

	// Ensure the bank got selected correctly.
	ptr, err := s.mbc1.BytePtr(lib.AccessTypeRead, 0x4000, 0)
	s.NoError(err)
	s.NotNil(ptr)
	s.Equal(uint8(0x7F), *ptr)

	ptr, err = s.mbc1.BytePtr(lib.AccessTypeRead, 0xA000, 0)
	s.NoError(err)
	s.NotNil(ptr)
	s.Equal(uint8(0x00), *ptr)

	// Now switch to RAM banking, which should reduce to bank 0x1F at most.
	_, err = s.mbc1.BytePtr(lib.AccessTypeWrite, 0x6000, 0x01)
	s.NoError(err)

	// Selected ROM bank is 0x1F!
	ptr, err = s.mbc1.BytePtr(lib.AccessTypeRead, 0x4000, 0)
	s.NoError(err)
	s.NotNil(ptr)
	s.Equal(uint8(0x1F), *ptr)

	// RAM selected bank is now 3.
	ptr, err = s.mbc1.BytePtr(lib.AccessTypeRead, 0xA000, 0)
	s.NoError(err)
	s.NotNil(ptr)
	s.Equal(uint8(0x03), *ptr)

	// Going back to ROM banking switches back the selector.
	_, err = s.mbc1.BytePtr(lib.AccessTypeWrite, 0x6000, 0x00)
	s.NoError(err)

	// Selected ROM bank is back to 0x7F.
	ptr, err = s.mbc1.BytePtr(lib.AccessTypeRead, 0x4000, 0)
	s.NoError(err)
	s.NotNil(ptr)
	s.Equal(uint8(0x7F), *ptr)

	// RAM selected bank is now 0.
	ptr, err = s.mbc1.BytePtr(lib.AccessTypeRead, 0xA000, 0)
	s.NoError(err)
	s.NotNil(ptr)
	s.Equal(uint8(0x00), *ptr)
}
