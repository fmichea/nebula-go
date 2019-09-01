package memory

import (
	"os"
)

var (
	ErrInvalidROMRead      = errors.New("failed to read ROM in full")
)

type MMU struct {
	ROM     []byte
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

	m.ROM = make([]byte, size)

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
	return nil
}
