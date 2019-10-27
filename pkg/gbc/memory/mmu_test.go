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

func (s *unitTestSuite) TestMMU_ReadByte_InMBC() {
	value, err := s.mmu.ReadByte(0x205)
	s.NoError(err)
	s.Equal(uint8(0x05), value)
}

func (s *unitTestSuite) TestMMU_ReadByte_InSegment() {
	value, err := s.mmu.ReadByte(0xC000)
	s.NoError(err)
	s.Equal(uint8(0), value)
}

func (s *unitTestSuite) TestMMU_ReadByte_InInvalidZone() {
	_, err := s.mmu.ReadByte(0xFEA0)
	s.Equal(lib.ErrInvalidRead, err)
}

func (s *unitTestSuite) TestMMU_ReadByte_InERAMWithNoERAM() {
	_, err := s.mmu.ReadByte(0xA000)
	s.Equal(lib.ErrInvalidRead, err)
}

func (s *unitTestSuite) TestMMU_ReadDByte_InMBC() {
	value, err := s.mmu.ReadDByte(0x205)
	s.NoError(err)
	s.Equal(uint16(0x0605), value)
}

func (s *unitTestSuite) TestMMU_ReadDByte_InERAMWithNoERAM_FirstByte() {
	_, err := s.mmu.ReadDByte(0xA000)
	s.Equal(lib.ErrInvalidRead, err)
}

func (s *unitTestSuite) TestMMU_ReadDByte_InERAMWithNoERAM_SecondByte() {
	_, err := s.mmu.ReadDByte(0xBFFF)
	s.Equal(lib.ErrInvalidRead, err)
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
