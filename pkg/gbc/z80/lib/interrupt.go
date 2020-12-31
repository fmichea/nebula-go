package z80lib

type Interrupt int

const (
	Rst00h Interrupt = iota
	Rst08h
	Rst10h
	Rst18h
	Rst20h
	Rst28h
	Rst30h
	Rst38h
	Rst40h
	Rst48h
	Rst50h
	Rst58h
	Rst60h
)

func (i Interrupt) Addr() uint16 {
	addrs := map[Interrupt]uint16{
		Rst00h: 0x0000,
		Rst08h: 0x0008,
		Rst10h: 0x0010,
		Rst18h: 0x0018,
		Rst20h: 0x0020,
		Rst28h: 0x0028,
		Rst30h: 0x0030,
		Rst38h: 0x0038,
		Rst40h: 0x0040,
		Rst48h: 0x0048,
		Rst50h: 0x0050,
		Rst58h: 0x0058,
		Rst60h: 0x0060,
	}
	return addrs[i]
}
