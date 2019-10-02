package mbcs

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"nebula-go/pkg/gbc/memory/lib"
	"nebula-go/pkg/gbc/memory/segments"
)

type unitTestSuite struct {
	suite.Suite

	rom  segments.Segment
	eram segments.Segment

	romOnly MBC
	mbc1    MBC
}

func generateData(bankSize, bankCount uint) []uint8 {
	totalSize := bankSize * bankCount

	data := make([]uint8, totalSize)
	for idx := uint(0); idx < totalSize; idx++ {
		data[idx] = uint8(idx / bankSize)
	}
	return data
}

func (s *unitTestSuite) SetupTest() {
	var err error

	romBankCount := lib.ROMSize2MB.BankCount()
	opts := []segments.Option{
		segments.WithBanks(romBankCount),
		segments.WithPinnedBank0(),
		segments.WithInitialData(generateData(0x4000, romBankCount)),
	}
	s.rom, err = segments.New(0x0000, 0x7FFF, opts...)
	s.NoError(err)

	ramBankCount := lib.RAMSize32KB.BankCount()
	opts = []segments.Option{
		segments.WithBanks(ramBankCount),
		segments.WithInitialData(generateData(0x2000, ramBankCount)),
	}
	s.eram, err = segments.New(0xA000, 0xBFFF, opts...)
	s.NoError(err)

	s.romOnly = newRomOnly(s.rom, s.eram)
	s.mbc1 = newMBC1(s.rom, s.eram)
}

func (s *unitTestSuite) testContainsAddress(mbc MBC) {
	s.True(mbc.ContainsAddress(0x0000))
	s.True(mbc.ContainsAddress(0x4000))
	s.True(mbc.ContainsAddress(0x7FFF))
	s.False(mbc.ContainsAddress(0x8000))

	s.True(mbc.ContainsAddress(0xA000))
	s.True(mbc.ContainsAddress(0xB000))
	s.True(mbc.ContainsAddress(0xBFFF))
	s.False(mbc.ContainsAddress(0xC000))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(unitTestSuite))
}
