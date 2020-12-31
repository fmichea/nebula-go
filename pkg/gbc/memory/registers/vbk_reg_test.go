package registers

import (
	"testing"

	"nebula-go/mocks/pkg/gbc/memory/segmentsmocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"nebula-go/pkg/common/testhelpers"
	"nebula-go/pkg/gbc/memory/cartridge"
	"nebula-go/pkg/gbc/memory/lib"
)

func TestVBKReg(t *testing.T) {
	runTest := func(t *testing.T, name string, fn func(mockVRAM *segmentsmocks.MockSegment)) {
		t.Run(name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockVRAM := segmentsmocks.NewMockSegment(mockCtrl)
			mockVRAM.EXPECT().BankCount().Return(uint(2)).AnyTimes()

			fn(mockVRAM)
		})
	}

	runTest(t, "DMG does not have any effect on VRAM", func(mockVRAM *segmentsmocks.MockSegment) {
		var vbkValue uint8

		vbk, err := NewVBKReg(&vbkValue, &cartridge.ROM{Type: lib.DMG01}, mockVRAM)
		require.NoError(t, err)

		value, err := vbk.Get()
		require.NoError(t, err)
		assert.Equal(t, uint8(0), value)

		require.NoError(t, vbk.Set(1))

		value, err = vbk.Get()
		require.NoError(t, err)
		assert.Equal(t, uint8(1), value)
	})

	runTest(t, "CGB select VRAM bank to 0 at init, changes set it", func(mockVRAM *segmentsmocks.MockSegment) {
		var vbkValue uint8

		mockVRAM.EXPECT().SelectBank(uint(0)).Return(nil)

		vbk, err := NewVBKReg(&vbkValue, &cartridge.ROM{Type: lib.CGB001}, mockVRAM)
		require.NoError(t, err)

		value, err := vbk.Get()
		require.NoError(t, err)
		assert.Equal(t, uint8(0), value)

		mockVRAM.EXPECT().SelectBank(uint(1)).Return(nil)
		require.NoError(t, vbk.Set(1))

		value, err = vbk.Get()
		require.NoError(t, err)
		assert.Equal(t, uint8(1), value)
	})

	runTest(t, "CGB with broken bank select errors out", func(mockVRAM *segmentsmocks.MockSegment) {
		var vbkValue uint8

		mockVRAM.EXPECT().SelectBank(uint(0)).Return(testhelpers.ErrTesting1)

		_, err := NewVBKReg(&vbkValue, &cartridge.ROM{Type: lib.CGB001}, mockVRAM)
		require.Equal(t, testhelpers.ErrTesting1, err)
	})
}
