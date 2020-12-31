package memory

import (
	"nebula-go/pkg/gbc/memory/mbcs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"nebula-go/pkg/gbc/memory/cartridge"
	"nebula-go/pkg/gbc/memory/lib"
)

type unitTestSuite struct {
	suite.Suite

	mmu *mmu
}

func (s *unitTestSuite) SetupTest() {
	var err error

	data := make([]uint8, 0x8000)
	for idx, value := range cartridge.NintendoLogo {
		data[0x104+idx] = value
	}
	data[0x14D] = 0xE7
	for idx := 0; idx < 0x100; idx++ {
		data[0x200+idx] = uint8(idx)
	}

	s.mmu, err = newMMUFromCartridge(&cartridge.ROM{
		Title:       "MMU TEST",
		Type:        lib.DMG01,
		Size:        lib.ROMSize32KB,
		Market:      lib.Japanese,
		RAMSize:     lib.RAMSizeNone,
		MBCSelector: mbcs.NewSelector(0),
		Data:        data,
	})
	s.NoError(err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(unitTestSuite))
}

func TestNewMMU(t *testing.T) {
	t.Run("invalid file refuses to create MMU", func(t *testing.T) {
		_, err := NewMMUFromFile(os.Stderr, "wedwedwedwed")
		assert.Error(t, err)
	})
}
