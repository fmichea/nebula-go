package opcodes

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"nebula-go/mocks/pkg/gbc/memorymocks"

	"nebula-go/pkg/gbc/z80/registers"
)

type unitTestSuite struct {
	suite.Suite

	mockCtrl *gomock.Controller
	mockMMU  *memorymocks.MockMMU

	regs    *registers.Registers
	factory *Factory
}

func (s *unitTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())

	s.mockMMU = memorymocks.NewMockMMU(s.mockCtrl)
	s.regs = registers.New()

	s.factory = NewFactory(s.mockMMU, s.regs)
}

func (s *unitTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *unitTestSuite) EqualFlags(expected uint8) {
	fmt := "flags were expected to be equal to \"%s\" but are \"%s\""

	flagsExpected := registers.NewFlags(expected)
	s.Equal(flagsExpected.Get(), s.regs.F.Get(), fmt, flagsExpected.String(), s.regs.F.String())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(unitTestSuite))
}
