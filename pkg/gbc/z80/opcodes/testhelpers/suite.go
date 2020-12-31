package testhelpers

import (
	"nebula-go/mocks/pkg/gbc/memorymocks"

	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/z80/registers"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type SetupFactoryFunc = func(mmu memory.MMU, regs *registers.Registers)

type OpcodesUnitTestSuite interface {
	suite.TestingSuite

	InitializeFactorySetup(fn SetupFactoryFunc)
	SetupTestFactory(mmu memory.MMU, regs *registers.Registers)
}

func Run(t *testing.T, s OpcodesUnitTestSuite) {
	s.InitializeFactorySetup(s.SetupTestFactory)
	suite.Run(t, s)
}

type OpcodesUnitTestSuiteMeta struct {
	suite.Suite

	mockCtrl *gomock.Controller

	MockMMU *memorymocks.MockMMU
	Regs    *registers.Registers

	setupFactoryFunc SetupFactoryFunc
}

func (s *OpcodesUnitTestSuiteMeta) InitializeFactorySetup(fn SetupFactoryFunc) {
	s.setupFactoryFunc = fn
}

func (s *OpcodesUnitTestSuiteMeta) SetupTest() {
	s.Require().NotNil(s.setupFactoryFunc)

	s.mockCtrl = gomock.NewController(s.T())

	s.MockMMU = memorymocks.NewMockMMU(s.mockCtrl)
	s.Regs = registers.New()

	s.setupFactoryFunc(s.MockMMU, s.Regs)
}

func (s *OpcodesUnitTestSuiteMeta) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *OpcodesUnitTestSuiteMeta) EqualFlags(expected uint8) {
	fmt := "flags were expected to be equal to \"%s\" but are \"%s\""

	flagsExpected := registers.NewFlags(expected)
	s.Equal(flagsExpected.Get(), s.Regs.F.Get(), fmt, flagsExpected.String(), s.Regs.F.String())
}
