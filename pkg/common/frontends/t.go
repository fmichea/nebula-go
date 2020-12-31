package frontends

import (
	"context"
)

type KeyboardState int
type KeyboardKey int

const (
	UpState KeyboardState = iota
	DownState
)

const (
	UpKey KeyboardKey = iota
	DownKey
	LeftKey
	RightKey

	AKey
	ZKey

	SpaceKey
	ReturnKey

	EscapeKey
)

type Pixel struct {
	r uint8
	g uint8
	b uint8

	transparent bool
}

func NewPixel(r, g, b uint8) Pixel {
	return Pixel{
		r: r,
		g: g,
		b: b,

		transparent: false,
	}
}

var (
	NonePixel        = Pixel{}
	TransparentPixel = Pixel{transparent: true}

	White = NewPixel(0xFF, 0xFF, 0xFF)
	Black = NewPixel(0x00, 0x00, 0x00)
)

type KeyboardCallback func(key KeyboardKey, state KeyboardState)

type MainWindow interface {
	Close() error

	SubscribeKeyboardStateChanges(callback KeyboardCallback)

	DrawLine(line uint, pixels []Pixel) error
	Commit() error

	MainLoop(ctx context.Context, cancel context.CancelFunc) error
}
