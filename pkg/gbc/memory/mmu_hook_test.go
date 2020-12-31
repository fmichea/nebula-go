package memory

import (
	"nebula-go/pkg/common/testhelpers"
	"nebula-go/pkg/gbc/memory/lib"
	"nebula-go/pkg/gbc/memory/mbcs"
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

func (s *unitTestSuite) TestMMU_ByteHook_OK() {
	var reg registerslib.Byte

	called := 0

	err := s.mmu.ByteHook(0xFFFE, func(ptr *uint8) (lib.Hook, error) {
		called++

		s.NotNil(ptr)

		reg = registerslib.NewByte(ptr, 0)
		return registerslib.WrapWithError(reg), nil
	})
	s.Require().NoError(err)
	s.Equal(1, called)

	s.Require().NotNil(reg)
	s.Equal(uint8(0), reg.Get())

	s.Require().NoError(s.mmu.WriteByte(0xFFFE, 0x10))
	s.Equal(uint8(0x10), reg.Get())
}

func (s *unitTestSuite) TestMMU_ByteHook_CannotHookByteTwice() {
	err := s.mmu.ByteHook(0xFFFE, func(ptr *uint8) (lib.Hook, error) {
		reg := registerslib.NewByte(ptr, 0)
		return registerslib.WrapWithError(reg), nil
	})
	s.Require().NoError(err)

	err = s.mmu.ByteHook(0xFFFE, nil)
	s.Require().Equal(lib.ErrDoubleHook, err)
}

func (s *unitTestSuite) TestMMU_ByteHook_CannotHookMBC() {
	err := s.mmu.ByteHook(0x100, nil)
	s.Equal(mbcs.ErrMBCHookInvalid, err)
}

func (s *unitTestSuite) TestMMU_ByteHook_CallbackReturnsError() {
	err := s.mmu.ByteHook(0xFFFE, func(ptr *uint8) (lib.Hook, error) {
		return nil, testhelpers.ErrTesting1
	})
	s.Equal(testhelpers.ErrTesting1, err)
}

func (s *unitTestSuite) TestMMU_ByteHook_CallbackDoesNotProvideHook() {
	err := s.mmu.ByteHook(0xFFFE, func(ptr *uint8) (lib.Hook, error) {
		return nil, nil
	})
	s.Equal(lib.ErrHookNotProvided, err)
}
