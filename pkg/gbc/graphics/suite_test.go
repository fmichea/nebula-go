package graphics

import (
	"nebula-go/mocks/pkg/gbc/memory/registersmocks"

	"nebula-go/pkg/gbc/memory/lib"
	"runtime/debug"
	"testing"
	"time"

	"nebula-go/mocks/pkg/common/clockmocks"
	"nebula-go/mocks/pkg/common/frontendsmocks"
	"nebula-go/mocks/pkg/gbc/memory/segmentsmocks"
	"nebula-go/mocks/pkg/gbc/memorymocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/suite"

	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/memory/cartridge"
	"nebula-go/pkg/gbc/memory/registers"
)

type unitTestSuiteTestContext struct {
	mockCtrl   *gomock.Controller
	mockMMU    *memorymocks.MockMMU
	mockVRAM   *segmentsmocks.MockSegment
	mockWRAM   *segmentsmocks.MockSegment
	mockWindow *frontendsmocks.MockMainWindow
	mockClock  *clockmocks.MockClock
	mockHDMA5  *registersmocks.MockHDMA5Reg

	originalTime time.Time

	mmuRegs *memory.Registers
	cr      *cartridge.ROM

	gpu *gpu
}

func initializePalette(index *registers.CGBPaletteIndexReg, reg *registers.CGBPaletteReg, value uint8) {
	index.AutoIncrement.SetBool(true)
	defer index.AutoIncrement.SetBool(false)

	for {
		reg.Set(value)

		if index.Index.Get() == 0 {
			break
		}
	}
}

func newUnitTestSuiteTestContext(t *testing.T) *unitTestSuiteTestContext {
	mockCtrl := gomock.NewController(t)

	uint8ptr := func() *uint8 {
		var val uint8
		return &val
	}

	obpi := registers.NewCGBPaletteIndexReg(uint8ptr())
	obpd := registers.NewCGBPaletteReg(obpi)

	bgpi := registers.NewCGBPaletteIndexReg(uint8ptr())
	bgpd := registers.NewCGBPaletteReg(bgpi)

	initializePalette(obpi, obpd, 0xFF)
	initializePalette(bgpi, bgpd, 0x00)

	cr := &cartridge.ROM{}

	mockVRAM := segmentsmocks.NewMockSegment(mockCtrl)
	mockVRAM.EXPECT().BankCount().Return(uint(1)).AnyTimes()

	mockWRAM := segmentsmocks.NewMockSegment(mockCtrl)
	mockWRAM.EXPECT().BankCount().Return(uint(1)).AnyTimes()

	mockMMU := memorymocks.NewMockMMU(mockCtrl)
	mockMMU.EXPECT().ByteHook(gomock.Any(), gomock.Any()).DoAndReturn(func(addr uint16, fn func(ptr *uint8) (hook lib.Hook, err error)) error {
		var value uint8

		_, err := fn(&value)
		return err
	}).AnyTimes()

	mmuRegs, err := memory.InitializeRegisters(mockMMU, cr, mockVRAM, mockWRAM)
	require.NoError(t, err)

	mockHDMA5 := registersmocks.NewMockHDMA5Reg(mockCtrl)
	mmuRegs.HDMA5 = mockHDMA5

	mockMMU.EXPECT().Registers().Return(mmuRegs).AnyTimes()
	mockMMU.EXPECT().Cartridge().Return(cr).AnyTimes()

	mockWindow := frontendsmocks.NewMockMainWindow(mockCtrl)

	originalTime := time.Now()

	mockClock := clockmocks.NewMockClock(mockCtrl)
	mockClock.EXPECT().Now().Return(originalTime)

	return &unitTestSuiteTestContext{
		mockCtrl:   mockCtrl,
		mockMMU:    mockMMU,
		mockVRAM:   mockVRAM,
		mockWRAM:   mockWRAM,
		mockWindow: mockWindow,
		mockHDMA5:  mockHDMA5,
		mockClock:  mockClock,

		originalTime: originalTime,

		mmuRegs: mmuRegs,
		cr:      cr,

		gpu: newGPU(mockMMU, mockWindow, NewPacer(mockClock)),
	}
}

func (c *unitTestSuiteTestContext) Finish() {
	c.mockCtrl.Finish()
}

type unitTestSuite struct {
	suite.Suite

	tctx *unitTestSuiteTestContext
}

func (s *unitTestSuite) SetupTest() {
	s.tctx = newUnitTestSuiteTestContext(s.T())
}

func (s *unitTestSuite) TearDownTest() {
	failOnPanic(s.T())

	s.tctx.Finish()
	s.tctx = nil
}

func failOnPanic(t *testing.T) {
	r := recover()
	if r != nil {
		t.Errorf("test panicked: %v\n%s", r, debug.Stack())
		t.FailNow()
	}
}

func (s *unitTestSuite) Run(name string, fn func()) {
	s.Suite.Run(name, func() {
		t := s.T()

		origTctx := s.tctx

		tctx := newUnitTestSuiteTestContext(t)
		defer func() {
			s.tctx = origTctx
			tctx.Finish()
		}()

		defer failOnPanic(t)

		s.tctx = tctx

		fn()
	})
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(unitTestSuite))
}
