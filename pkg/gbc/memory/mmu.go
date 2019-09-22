package memory

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"nebula-go/pkg/gbc/memory/lib"
)

var (
	ErrInvalidROMRead      = errors.New("failed to read ROM in full")
	ErrNintendoLogoInvalid = errors.New("nintendo logo from ROM did not match")
	ErrMBCNotImplemented   = errors.New("memory controller has not been implemented yet")
	ErrROMSizeInvalid      = errors.New("got an invalid ROM size value")
	ErrRAMSizeInvalid      = errors.New("got an invalid RAM size value")
	ErrChecksumInvalid     = errors.New("cartridge checksum is invalid")
)

const (
	CGBFlagAddress             = 0x143
	CGBFlagBackwardsCompatible = 0x80
	CGBFlagColorOnly           = 0xC0
)

var (
	NintendoLogo = []uint8(
		"\xce\xed\x66\x66\xcc\x0d\x00\x0b\x03\x73\x00\x83\x00\x0c\x00\x0d" +
			"\x00\x08\x11\x1f\x88\x89\x00\x0e\xdc\xcc\x6e\xe6\xdd\xdd\xd9\x99" +
			"\xbb\xbb\x67\x63\x6e\x0e\xec\xcc\xdd\xdc\x99\x9f\xbb\xb9\x33\x3e")

	NintendoLogoAddress    = 0x104
	NintendoLogoEndAddress = NintendoLogoAddress + len(NintendoLogo)

	TitleAddress  = 0x134
	TargetAddress = 0x14A

	ROMSizeAddress = 0x148
	RAMSizeAddress = 0x149

	MBCFlagAddress = 0x0147
	MBCNames       = []string{
		/* 0x00 */
		"ROM ONLY",
		/* 0x01 -> 0x03: MBC1 */
		"MBC1",
		"MBC1+RAM",
		"MBC1+RAM+BATTERY",
		/* 0x04: Unused */
		"",
		/* 0x05 -> 0x06: MBC2 */
		"MBC2",
		"MBC2+BATTERY",
		/* 0x07: Unused */
		"",
		/* 0x08 -> 0x09: ROM+RAM */
		"ROM+RAM",
		"ROM+RAM+BATTERY",
		/* 0x0A: Unused */
		"",
		/* 0x0B -> 0x0D: MMM01 */
		"MMM01",
		"MMM01+RAM",
		"MMM01+RAM+BATTERY",
		/* 0x0E: Unused */
		"",
		/* 0x0F -> 0x13: MBC3 */
		"MBC3+TIMER+BATTERY",
		"MBC3+TIMER+RAM+BATTERY",
		"MBC3",
		"MBC3+RAM",
		"MBC3+RAM+BATTERY",
		/* 0x14: Unused */
		"",
		/* 0x15 -> 0x17: MBC4 */
		"MBC4",
		"MBC4+RAM",
		"MBC4+RAM+BATTERY",
		/* 0x18: Unused */
		"",
		/* 0x19 -> 0x1E: MBC5 */
		"MBC5",
		"MBC5+RAM",
		"MBC5+RAM+BATTERY",
		"MBC5+RUMBLE+RAM",
		"MBC5+RUMBLE+RAM+BATTERY",
	}

	MBCNamesCount = uint8(len(MBCNames))

	ChecksumStartAddress = 0x134
	ChecksumEndAddress   = 0x14D
)

type MMU struct {
	Title string

	ROMType lib.ROMType
	ROM     []uint8
}

func NewMMU(filename string) (*MMU, error) {
	result := &MMU{}

	if err := result.loadRom(filename); err != nil {
		return nil, err
	}

	if err := result.checkROM(); err != nil {
		return nil, err
	}

	return result, nil
}

func (m *MMU) loadRom(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	size := stat.Size()

	m.ROM = make([]uint8, size)

	readCount, err := f.Read(m.ROM)
	if err != nil {
		return err
	}

	if readCount != int(size) {
		return ErrInvalidROMRead
	}

	return nil
}

func (m *MMU) checkROM() error {
	nintendoLogo := m.ROM[NintendoLogoAddress:NintendoLogoEndAddress]
	if !bytes.Equal(NintendoLogo, nintendoLogo) {
		return ErrNintendoLogoInvalid
	}

	value := m.ROM[CGBFlagAddress]
	if value == CGBFlagBackwardsCompatible || value == CGBFlagColorOnly {
		m.ROMType = lib.CGB001
	} else {
		m.ROMType = lib.DMG01
	}

	fmt.Println("ROM Type:", m.ROMType)

	titleBytes := m.ROM[TitleAddress : TitleAddress+m.ROMType.GetTitleSize()]
	zeroByteIndex := bytes.IndexByte(titleBytes, 0)
	m.Title = string(titleBytes[:zeroByteIndex])

	fmt.Println("ROM Title:", m.Title)

	if err := m.loadMBC(); err != nil {
		return err
	}

	target := Japanese
	if m.ROM[TargetAddress] != 0 {
		target = NonJapanese
	}
	fmt.Println("ROM Target:", target)

	if err := m.loadROMSize(); err != nil {
		return err
	}

	if err := m.loadRAMSize(); err != nil {
		return err
	}

	if err := m.verifyHeaderChecksum(); err != nil {
		return err
	} else {
		fmt.Println("Cartridge checksum valid!")
	}

	return nil
}

func (m *MMU) verifyHeaderChecksum() error {
	checksum := uint8(0)

	addr := ChecksumStartAddress
	for addr <= ChecksumEndAddress {
		checksum -= m.ROM[addr] + 1
		addr++
	}

	if checksum != 0xFF {
		return ErrChecksumInvalid
	}

	return nil
}

func (m *MMU) getMBCFromSelectorValue(value uint8) MBC {
	switch value {
	case 0x00:
	// load ROMonly MBC

	case 0x01, 0x02, 0x03:
		return &MBC1{}

	case 0x05, 0x06:
	// load MBC2

	case 0x0F, 0x10, 0x11, 0x12, 0x13:
	// load MBC3

	case 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E:
	// load MBC5

	default:
		// log bad thing?
	}
	return nil
}

func (m *MMU) getMBCNameFromSelectorValue(value uint8) string {
	if value < MBCNamesCount {
		return MBCNames[value]
	} else {
		switch value {
		case 0xFC:
			return "POCKET CAMERA"
		case 0xFD:
			return "BANDAI TAMA5"
		case 0xFE:
			return "HuC3"
		case 0xFF:
			return "HuC1+RAM+BATTERY"
		}
	}
	return ""
}

func (m *MMU) loadMBC() error {
	value := m.ROM[MBCFlagAddress]

	mbc := m.getMBCFromSelectorValue(value)
	name := m.getMBCNameFromSelectorValue(value)

	if name != "" {
		fmt.Println("Cartridge type:", name)
	} else {
		fmt.Printf("Cartridge type is not a known value: %02x\n", value)
	}

	if mbc == nil {
		return ErrMBCNotImplemented
	}

	return nil
}

func (m *MMU) loadROMSize() error {
	romSize := lib.ROMSize(m.ROM[ROMSizeAddress])

	if !romSize.IsValid() {
		return ErrROMSizeInvalid
	}

	fmt.Println("ROM Size:", romSize)
	return nil
}

func (m *MMU) loadRAMSize() error {
	ramSize := lib.RAMSize(m.ROM[RAMSizeAddress])
	if !ramSize.IsValid() {
		return ErrRAMSizeInvalid
	}
	fmt.Println("RAM Size:", ramSize)
	return nil
}

func (m *MMU) ReadByte(addr uint16) uint8 {
	return m.readByteInternal(addr)
}

func (m *MMU) ReadDByte(addr uint16) uint16 {
	result := uint16(m.readByteInternal(addr+1)) << 8
	result |= uint16(m.readByteInternal(addr))
	return result
}

func (m *MMU) readByteInternal(addr uint16) uint8 {
	return 0
}
