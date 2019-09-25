package cartridge

import (
	"fmt"
	"io"

	"nebula-go/pkg/gbc/memory/lib"
	"nebula-go/pkg/gbc/memory/mbcs"
)

type ROM struct {
	Title string

	Type   lib.ROMType
	Size   lib.ROMSize
	Market lib.ROMMarket

	RAMSize lib.RAMSize

	MBCSelector *mbcs.Selector

	Data []uint8
}

func (r *ROM) PrintInformation(out io.Writer) error {
	if _, err := fmt.Fprintln(out, "===== ROM INFORMATION ====="); err != nil {
		return err
	}

	for _, info := range []struct {
		name  string
		value string
	}{
		{"Device", r.Type.String()},
		{"ROM Title", r.Title},
		{"ROM Market", r.Market.String()},
		{"MBC Controller", r.MBCSelector.Name()},
		{"ROM Size", r.Size.String()},
		{"RAM Size", r.RAMSize.String()},
	} {
		if _, err := fmt.Fprintf(out, "%-20s%s\n", fmt.Sprintf("%s:", info.name), info.value); err != nil {
			return err
		}
	}

	return nil
}
