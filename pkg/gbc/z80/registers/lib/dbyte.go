package registerslib

type DByte interface {
	Get() uint16
	Set(value uint16)
}

type dbyte struct {
	value uint16
}

func NewDByte(value uint16) DByte {
	return &dbyte{
		value: value,
	}
}

func (d *dbyte) Get() uint16 {
	return d.value
}

func (d *dbyte) Set(value uint16) {
	d.value = value
}
