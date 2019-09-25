package cartridge

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Run("file read fails", func(t *testing.T) {
		var out bytes.Buffer

		rom, err := Load(&out, "wedwedwed")
		require.Nil(t, rom)
		require.Error(t, err)
	})

	t.Run("not a cartridge", func(t *testing.T) {
		tmpFile, err := ioutil.TempFile(os.TempDir(), "nebula-go-")
		require.NoError(t, err)

		defer func() {
			require.NoError(t, os.Remove(tmpFile.Name()))
		}()

		require.NoError(t, tmpFile.Close())

		var out bytes.Buffer

		rom, err := Load(&out, tmpFile.Name())
		require.Nil(t, rom)
		require.Error(t, err)
	})

	t.Run("valid case", func(t *testing.T) {
		var out bytes.Buffer

		data := make([]uint8, _minimumROMDataSize)
		for idx, value := range _nintendoLogo {
			data[_nintendoLogoStartAddress+idx] = value
		}

		checksumValue := uint8(_checksumEndAddress - _checksumStartAddress + 1)
		for idx, value := range []uint8("TEST GAME") {
			checksumValue += value
			data[_titleAddress+idx] = value
		}
		data[_checksumEndAddress] = (0xFF - checksumValue) + 2

		tmpFile, err := ioutil.TempFile(os.TempDir(), "nebula-go-")
		require.NoError(t, err)

		defer func() {
			require.NoError(t, os.Remove(tmpFile.Name()))
		}()

		_, err = tmpFile.Write(data)
		require.NoError(t, err)
		require.NoError(t, tmpFile.Close())

		rom, err := Load(&out, tmpFile.Name())
		require.NoError(t, err)
		require.NotNil(t, rom)

		for _, line := range []string{
			"Device:             GAME BOY (DMG-01)",
			"ROM Title:          TEST GAME",
			"ROM Market:         Japanese",
			"MBC Controller:     ROM ONLY",
			"ROM Size:           32KByte",
			"RAM Size:           None",
		} {
			assert.Contains(t, out.String(), line)
		}
	})
}
