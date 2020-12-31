package z80

import (
	"nebula-go/pkg/common/testhelpers"
)

func (s *unitTestSuite) TestOpcodeInitialization() {
	s.Len(s.cpu.Opcodes, 0x100)
}

func (s *unitTestSuite) TestDoCycle() {
	s.mockMMU.EXPECT().ReadByte(uint16(0x100)).Return(uint8(0x00), nil)

	s.mockMMU.EXPECT().ReadByte(uint16(0x101)).Return(uint8(0x00), nil)
	s.mockMMU.EXPECT().ReadByte(uint16(0x102)).Return(uint8(0x00), nil)

	s.mockGPU.EXPECT().DoCycles(uint16(4)).Return(nil)

	s.Require().NoError(s.cpu.doCycle(true))
}

func (s *unitTestSuite) TestDoCycle_ErrorReadingPC() {
	s.mockMMU.EXPECT().ReadByte(uint16(0x100)).Return(uint8(0x00), testhelpers.ErrTesting1)

	s.Require().Equal(testhelpers.ErrTesting1, s.cpu.doCycle(false))
}

func (s *unitTestSuite) TestDoCycle_ErrorExecutingOpcode() {
	s.mockMMU.EXPECT().ReadByte(uint16(0x100)).Return(uint8(0x08), nil)
	s.mockMMU.EXPECT().ReadDByte(uint16(0x101)).Return(uint16(0), testhelpers.ErrTesting1)

	s.Require().Equal(testhelpers.ErrTesting1, s.cpu.doCycle(false))
}

func (s *unitTestSuite) TestDoCycle_ErrorDoingGPUCycle() {
	s.mockMMU.EXPECT().ReadByte(uint16(0x100)).Return(uint8(0x00), nil)
	s.mockGPU.EXPECT().DoCycles(uint16(4)).Return(testhelpers.ErrTesting1)

	s.Require().Equal(testhelpers.ErrTesting1, s.cpu.doCycle(false))
}

func (s *unitTestSuite) TestDoCycle_DoubleSpeedMakeGPUGoTwiceAsSlow() {
	s.mockRegs.KEY1.ChangeRequest.SetBool(true)
	s.mockRegs.KEY1.SwitchIfRequested()

	s.mockMMU.EXPECT().ReadByte(uint16(0x100)).Return(uint8(0x00), nil)
	s.mockGPU.EXPECT().DoCycles(uint16(2)).Return(nil)

	s.Require().NoError(s.cpu.doCycle(false))
}

func (s *unitTestSuite) TestRun_ExecuteThreeNoopThenError() {
	s.mockMMU.EXPECT().ReadByte(uint16(0x100)).Return(uint8(0x00), nil)
	s.mockMMU.EXPECT().ReadByte(uint16(0x101)).Return(uint8(0x00), nil)
	s.mockMMU.EXPECT().ReadByte(uint16(0x102)).Return(uint8(0x00), nil)
	s.mockMMU.EXPECT().ReadByte(uint16(0x103)).Return(uint8(0x00), testhelpers.ErrTesting1)

	s.mockGPU.EXPECT().DoCycles(uint16(4)).Return(nil).Times(3)

	s.Require().Equal(testhelpers.ErrTesting1, s.cpu.Run())
}
