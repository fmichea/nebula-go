package joypad

import (
	"testing"

	"nebula-go/mocks/pkg/common/frontendsmocks"
	"nebula-go/mocks/pkg/gbc/memorymocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"nebula-go/pkg/common/frontends"
	"nebula-go/pkg/common/ptr"
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/memory/registers"
)

type unitTestSuite struct {
	suite.Suite

	mockCtrl *gomock.Controller

	mockRegs   *memory.Registers
	mockMMU    *memorymocks.MockMMU
	mockWindow *frontendsmocks.MockMainWindow

	callback frontends.KeyboardCallback

	joypad Joypad
}

func (s *unitTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())

	s.mockMMU = memorymocks.NewMockMMU(s.mockCtrl)

	s.mockRegs = &memory.Registers{
		IF:   registers.NewInterruptReg(ptr.UInt8(0)),
		JOYP: registers.NewJOYPReg(ptr.UInt8(0)),
	}
	s.mockMMU.EXPECT().Registers().Return(s.mockRegs).AnyTimes()

	s.mockWindow = frontendsmocks.NewMockMainWindow(s.mockCtrl)
	s.mockWindow.EXPECT().SubscribeKeyboardStateChanges(gomock.Any()).Do(func(callback frontends.KeyboardCallback) {
		s.callback = callback
	})

	s.joypad = New(s.mockMMU, s.mockWindow)
}

func (s *unitTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(unitTestSuite))
}
