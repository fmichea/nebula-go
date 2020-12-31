package memory

import (
	"nebula-go/pkg/gbc/memory/lib"
)

func (m *mmu) ByteHook(addr uint16, fn func(ptr *uint8) (lib.Hook, error)) error {
	if _, ok := m.hooks[addr]; ok {
		return lib.ErrDoubleHook
	}

	segment, err := m.getSegmentForAddress(addr)
	if err != nil {
		return err
	}

	ptr, err := segment.ByteHook(addr)
	if err != nil {
		return err
	}

	hook, err := fn(ptr)
	if err != nil {
		return err
	}

	if hook == nil {
		return lib.ErrHookNotProvided
	}

	m.hooks[addr] = hook
	return nil
}
