package cb

import (
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/z80/opcodes/testhelpers"
	"nebula-go/pkg/gbc/z80/registers"

	"testing"
)

type cbBitOpTestCase struct {
	initialValue uint8
	bit          uint8
	resultValue  uint8

	initialFlags uint8
	resultFlags  uint8
}

type unitTestSuite struct {
	testhelpers.OpcodesUnitTestSuiteMeta

	factory *Factory
}

func (s *unitTestSuite) SetupTestFactory(mmu memory.MMU, regs *registers.Registers) {
	s.factory = NewFactory(mmu, regs)
}

func TestSuite(t *testing.T) {
	testhelpers.Run(t, &unitTestSuite{})
}
