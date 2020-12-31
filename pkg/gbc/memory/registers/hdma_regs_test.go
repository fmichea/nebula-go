package registers

import (
	"fmt"
	"testing"

	"nebula-go/mocks/pkg/gbc/memory/libmocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"nebula-go/pkg/common/testhelpers"
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"

	"github.com/golang/mock/gomock"
)

type hdmaTestContext struct {
	mockMMU *libmocks.MockMemoryIO

	hdma1 registerslib.Byte
	hdma2 registerslib.Byte
	hdma3 registerslib.Byte
	hdma4 registerslib.Byte
	hdma5 HDMA5Reg
}

func runHDMATest(t *testing.T, name string, fn func(t *testing.T, tctx hdmaTestContext)) {
	uint8ptr := func(value uint8) *uint8 {
		return &value
	}

	t.Run(name, func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockMMU := libmocks.NewMockMemoryIO(mockCtrl)

		hdma1 := NewHDMA1Reg(uint8ptr(0))
		hdma1.Set(0xAA)

		hdma2 := NewHDMA2Reg(uint8ptr(0))
		hdma2.Set(0xBB)

		hdma3 := NewHDMA3Reg(uint8ptr(0))
		hdma3.Set(0xCC)

		hdma4 := NewHDMA4Reg(uint8ptr(0))
		hdma4.Set(0xDD)

		fn(t, hdmaTestContext{
			mockMMU: mockMMU,

			hdma1: hdma1,
			hdma2: hdma2,
			hdma3: hdma3,
			hdma4: hdma4,

			hdma5: NewHDMA5Reg(uint8ptr(0xFF), mockMMU, hdma1, hdma2, hdma3, hdma4),
		})
	})
}

func TestHDMA5(t *testing.T) {
	for _, c := range []struct {
		value  uint8
		length uint
	}{
		{0x00, 0x10},
		{0x03, 0x40},
		{0x7F, 0x800},
	} {
		name := fmt.Sprintf("writing with bit7=0 copies length immediately: %v", c)

		runHDMATest(t, name, func(t *testing.T, tctx hdmaTestContext) {
			buffer := make([]uint8, c.length)

			tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xAAB0), c.length).Return(buffer, nil)
			tctx.mockMMU.EXPECT().WriteByteSlice(uint16(0x8CD0), buffer).Return(nil)

			require.NoError(t, tctx.hdma5.Set(c.value))

			value, err := tctx.hdma5.Get()
			require.NoError(t, err)
			assert.Equal(t, uint8(0xFF), value)
		})
	}

	runHDMATest(t, "writing with bit7=0 copy error: invalid read", func(t *testing.T, tctx hdmaTestContext) {
		tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xAAB0), uint(0x10)).Return(nil, testhelpers.ErrTesting1)

		err := tctx.hdma5.Set(0x00)
		assert.Equal(t, testhelpers.ErrTesting1, err)
	})

	runHDMATest(t, "writing with bit7=0 copy error: invalid write", func(t *testing.T, tctx hdmaTestContext) {
		buffer := make([]uint8, 0x10)

		tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xAAB0), uint(0x10)).Return(buffer, nil)
		tctx.mockMMU.EXPECT().WriteByteSlice(uint16(0x8CD0), buffer).Return(testhelpers.ErrTesting1)

		err := tctx.hdma5.Set(0x00)
		assert.Equal(t, testhelpers.ErrTesting1, err)
	})

	runHDMATest(t, "writing with bit7=1 activates copy in horizontal DMA", func(t *testing.T, tctx hdmaTestContext) {
		buffer := make([]uint8, 0x10)

		// This will not copy anything, but it will start HDMA.
		require.NoError(t, tctx.hdma5.Set(0x82))

		// The value returned is the length left to copy.
		value, err := tctx.hdma5.Get()
		require.NoError(t, err)
		assert.Equal(t, uint8(0x2), value)

		// Calling MaybeDoHDMA copies 0x10 bytes.
		tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xAAB0), uint(0x10)).Return(buffer, nil)
		tctx.mockMMU.EXPECT().WriteByteSlice(uint16(0x8CD0), buffer).Return(nil)

		require.NoError(t, tctx.hdma5.MaybeDoHDMA())

		// Length left to copy is decremented.
		value, err = tctx.hdma5.Get()
		require.NoError(t, err)
		assert.Equal(t, uint8(0x1), value)

		// Calling MaybeDoHDMA copies another 0x10 bytes.
		tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xAAC0), uint(0x10)).Return(buffer, nil)
		tctx.mockMMU.EXPECT().WriteByteSlice(uint16(0x8CE0), buffer).Return(nil)

		require.NoError(t, tctx.hdma5.MaybeDoHDMA())

		// Length left to copy is decremented.
		value, err = tctx.hdma5.Get()
		require.NoError(t, err)
		assert.Equal(t, uint8(0x0), value)

		// Another call to MaybeDoHDMA copies the last 0x10 bytes.
		tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xAAD0), uint(0x10)).Return(buffer, nil)
		tctx.mockMMU.EXPECT().WriteByteSlice(uint16(0x8CF0), buffer).Return(nil)

		require.NoError(t, tctx.hdma5.MaybeDoHDMA())

		// Length is reset to 0xFF now.
		value, err = tctx.hdma5.Get()
		require.NoError(t, err)
		assert.Equal(t, uint8(0xFF), value)

		// Calling MaybeDoHDMA doesn't do anything when bit7 is set.
		require.NoError(t, tctx.hdma5.MaybeDoHDMA())
	})

	runHDMATest(t, "writing with bit7=1 activates copy in horizontal DMA, error case", func(t *testing.T, tctx hdmaTestContext) {
		// This will not copy anything, but it will start HDMA.
		require.NoError(t, tctx.hdma5.Set(0x82))

		// Expect the read to fail.
		tctx.mockMMU.EXPECT().ReadByteSlice(uint16(0xAAB0), uint(0x10)).Return(nil, testhelpers.ErrTesting1)

		// Calling MaybeDoHDMA fails on read, error is returned.
		assert.Equal(t, testhelpers.ErrTesting1, tctx.hdma5.MaybeDoHDMA())
	})
}
