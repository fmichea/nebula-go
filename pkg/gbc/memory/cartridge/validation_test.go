package cartridge

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	outputBuilder := func(lines []string) string {
		result := "===== DATA VALIDATION =====\n"
		result += strings.Join(lines, "\n")
		result += "\n"
		return result
	}

	t.Run("ROM file data is too small (no space for headers)", func(t *testing.T) {
		var out bytes.Buffer

		err := validate(&out, nil)
		assert.Equal(t, ErrInvalidROMRead, err)

		output := outputBuilder([]string{
			"Checking ROM file data...               FAILED",
		})
		assert.Equal(t, output, out.String())
	})

	t.Run("valid case, all other check functions are tested separately", func(t *testing.T) {
		var out bytes.Buffer

		data := make([]uint8, _minimumROMDataSize)
		for idx, value := range NintendoLogo {
			data[_nintendoLogoStartAddress+idx] = value
		}
		data[_checksumEndAddress] = 0xE7 // empty rom checksum

		err := validate(&out, data)
		require.NoError(t, err)

		output := outputBuilder([]string{
			"Checking ROM file data...               PASSED",
			"Checking NINTENDO logo...               PASSED",
			"Checking cartridge header checksum...   PASSED",
			"Checking ROM size flag...               PASSED",
			"Checking RAM size flag...               PASSED",
			"Checking MBC selector flag...           PASSED",
		})
		assert.Equal(t, output, out.String())
	})
}
