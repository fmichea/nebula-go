package controlflow

import (
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/z80/opcodes/testhelpers"
	"nebula-go/pkg/gbc/z80/registers"

	"testing"
)

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
