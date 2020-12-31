package lib

type Hook interface {
	Get() (uint8, error)
	Set(value uint8) error
}
