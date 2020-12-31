package memory

func (s *unitTestSuite) TestMMU_WriteByte_CannotWriteMBC() {
	err := s.mmu.WriteByte(0x100, 0xFF)
	s.EqualError(err, "write error at 0100h: invalid write")
}

func (s *unitTestSuite) TestMMU_WriteDByte_FirstByteFailed() {
	err := s.mmu.WriteDByte(0x100, 0xFF)
	s.EqualError(err, "write error at 0100h: invalid write")
}

func (s *unitTestSuite) TestMMU_WriteDByte_SecondByteFailed() {
	err := s.mmu.WriteDByte(0xFFFF, 0xFFFF)
	s.EqualError(err, "write error at FFFFh: invalid write")
}
