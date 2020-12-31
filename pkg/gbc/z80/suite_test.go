package z80

import (
	"testing"

	"nebula-go/mocks/pkg/gbc/graphicsmocks"
	"nebula-go/mocks/pkg/gbc/memorymocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/memory/registers"
)

type unitTestSuite struct {
	suite.Suite

	mockCtrl *gomock.Controller

	mockRegs *memory.Registers
	mockMMU  *memorymocks.MockMMU
	mockGPU  *graphicsmocks.MockGPU

	cpu *CPU
}

func (s *unitTestSuite) SetupTest() {
	uint8ptr := func() *uint8 {
		var val uint8
		return &val
	}

	s.mockCtrl = gomock.NewController(s.T())

	s.mockMMU = memorymocks.NewMockMMU(s.mockCtrl)
	s.mockGPU = graphicsmocks.NewMockGPU(s.mockCtrl)

	tac := registers.NewTACReg(uint8ptr())
	div := registers.NewDIVReg(uint8ptr(), tac)

	s.mockRegs = &memory.Registers{
		KEY1: registers.NewKEY1Reg(uint8ptr()),
		DIV:  div,
		TIMA: registers.NewTIMAReg(uint8ptr(), div, nil, nil),

		IF: registers.NewInterruptReg(uint8ptr()),
		IE: registers.NewInterruptReg(uint8ptr()),
	}
	s.mockMMU.EXPECT().Registers().Return(s.mockRegs).AnyTimes()

	s.cpu = NewCPU(s.mockMMU, s.mockGPU)
}

func (s *unitTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(unitTestSuite))
}
