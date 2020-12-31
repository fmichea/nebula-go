package registerslib

type ByteWithErrors interface {
	Get() (uint8, error)
	Set(value uint8) error
}

func WrapWithError(b Byte) *byteWithError {
	return &byteWithError{reg: b}
}

type byteWithError struct {
	reg Byte
}

func (b *byteWithError) Get() (uint8, error) {
	return b.reg.Get(), nil
}

func (b *byteWithError) Set(value uint8) error {
	b.reg.Set(value)
	return nil
}
