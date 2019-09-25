package cartridge

import (
	"fmt"
	"io"

	"go.uber.org/multierr"
)

const _minimumROMDataSize = 0x14F

func verifyROMDataMinimumSize(romData []uint8) error {
	if len(romData) < _minimumROMDataSize {
		return ErrInvalidROMRead
	}
	return nil
}

type check struct {
	name string
	fn   func(romData []uint8) error
}

func validate(out io.Writer, romData []uint8) error {
	if _, err := fmt.Fprintln(out, "===== DATA VALIDATION ====="); err != nil {
		return err
	}

	for _, dv := range []check{
		{"ROM file data", verifyROMDataMinimumSize},
		{"NINTENDO logo", verifyNintendoLogo},
		{"cartridge header checksum", verifyHeaderChecksum},
		{"ROM size flag", verifyROMSizeFlag},
		{"RAM size flag", verifyRAMSizeFlag},
		{"MBC selector flag", verifyMBCSelector},
	} {
		if _, err := fmt.Fprintf(out, "%-40s", fmt.Sprintf("Checking %s... ", dv.name)); err != nil {
			return err
		}

		if err := dv.fn(romData); err != nil {
			if _, err2 := fmt.Fprintf(out, "FAILED\n"); err2 != nil {
				return multierr.Combine(err, err2)
			}
			return err
		}

		if _, err := fmt.Fprintf(out, "PASSED\n"); err != nil {
			return err
		}
	}
	return nil
}
