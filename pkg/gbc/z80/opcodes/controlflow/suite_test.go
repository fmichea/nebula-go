package controlflow

import (
	"nebula-go/pkg/gbc/memory"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	"nebula-go/pkg/gbc/z80/opcodes/testhelpers"
	"testing"
)

type unitTestSuite struct {
	testhelpers.OpcodesUnitTestSuiteMeta

	factory *Factory
}

func (s *unitTestSuite) SetupTestFactory(mmu memory.MMU, regs *z80lib.Registers) {
	s.factory = NewFactory(mmu, regs)
}

func TestSuite(t *testing.T) {
	testhelpers.Run(t, &unitTestSuite{})
}
