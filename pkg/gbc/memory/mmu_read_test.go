package memory

func (s *unitTestSuite) TestMMU_ReadByte_InMBC() {
	value, err := s.mmu.ReadByte(0x205)
	s.NoError(err)
	s.Equal(uint8(0x05), value)
}

func (s *unitTestSuite) TestMMU_ReadByte_InSegment() {
	value, err := s.mmu.ReadByte(0xC000)
	s.NoError(err)
	s.Equal(uint8(0), value)
}

func (s *unitTestSuite) TestMMU_ReadByte_InInvalidZone() {
	_, err := s.mmu.ReadByte(0xFEA0)
	s.EqualError(err, "read error at FEA0h: no segment at address")
}

func (s *unitTestSuite) TestMMU_ReadByte_InERAMWithNoERAM() {
	_, err := s.mmu.ReadByte(0xA000)
	s.EqualError(err, "read error at A000h: invalid read")
}

func (s *unitTestSuite) TestMMU_ReadDByte_InMBC() {
	value, err := s.mmu.ReadDByte(0x205)
	s.NoError(err)
	s.Equal(uint16(0x0605), value)
}

func (s *unitTestSuite) TestMMU_ReadDByte_InERAMWithNoERAM_FirstByte() {
	_, err := s.mmu.ReadDByte(0xA000)
	s.EqualError(err, "read error at A000h: invalid read")
}

func (s *unitTestSuite) TestMMU_ReadDByte_InERAMWithNoERAM_SecondByte() {
	_, err := s.mmu.ReadDByte(0xBFFF)
	s.EqualError(err, "read error at BFFFh: invalid read")
}
