package cartridge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyNintendoLogo(t *testing.T) {
	setup := func(copyLogo bool) []uint8 {
		romData := make([]uint8, _minimumROMDataSize)
		if copyLogo {
			for idx, value := range _nintendoLogo {
				romData[_nintendoLogoStartAddress+idx] = value
			}
		}
		return romData
	}

	t.Run("invalid data does not pass check", func(t *testing.T) {
		assert.Equal(t, ErrNintendoLogoInvalid, verifyNintendoLogo(setup(false)))
	})

	t.Run("valid data passes check", func(t *testing.T) {
		assert.NoError(t, verifyNintendoLogo(setup(true)))
	})
}
