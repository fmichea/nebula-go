package cartridge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyHeaderChecksum(t *testing.T) {
	setup := func() []uint8 {
		return make([]uint8, _minimumROMDataSize)
	}

	t.Run("invalid ROM does not pass checksum", func(t *testing.T) {
		data := setup()
		assert.Equal(t, ErrChecksumInvalid, verifyHeaderChecksum(data))
	})

	t.Run("valid ROM data passing the chechsum", func(t *testing.T) {
		data := setup()

		// As a sum, all of the ones added, plus the values:
		//     (0x14D - 0x134 + 1) + 0x80 + 0x38 + 0x46 + 0xe9 = 0x201
		// Therefore, subtracted to a byte, it would result in 0xFF, valid checksum.
		for idx, value := range []uint8("\x80\x38\x46\xe9") {
			data[_checksumStartAddress+idx] = value
		}
		assert.NoError(t, verifyHeaderChecksum(data))
	})
}
