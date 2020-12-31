package frontends

import (
	"context"
	"errors"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	"go.uber.org/multierr"
)

const (
	_defaultZoom = 3
	_pixelSize   = 4

	_sdlSubSystemFlags = sdl.INIT_VIDEO | sdl.INIT_JOYSTICK | sdl.INIT_EVENTS

	_sdlWindowPos = sdl.WINDOWPOS_UNDEFINED

	_sdlRendererFlags = sdl.RENDERER_ACCELERATED | sdl.RENDERER_PRESENTVSYNC

	_sdlTextureFormat = sdl.PIXELFORMAT_ARGB8888
	_sdlTextureAccess = sdl.TEXTUREACCESS_STREAMING
)

var (
	_sdlWindowFlags = sdl.WINDOW_SHOWN | sdl.WINDOW_RESIZABLE

	_sdlKeys = map[sdl.Keycode]KeyboardKey{
		sdl.K_UP:    UpKey,
		sdl.K_DOWN:  DownKey,
		sdl.K_LEFT:  LeftKey,
		sdl.K_RIGHT: RightKey,

		sdl.K_a: AKey,
		sdl.K_z: ZKey,

		sdl.K_SPACE:  SpaceKey,
		sdl.K_RETURN: ReturnKey,

		sdl.K_ESCAPE: EscapeKey,
	}
)

func sdlData(p Pixel) []byte {
	return []byte{p.b, p.g, p.r, 0x00}
}

func windowSize(width, height, zoom int32) (int32, int32) {
	return width * zoom, height * zoom
}

func rendererPitch(width, zoom int32) int {
	return _pixelSize * int(width*zoom)
}

type mainWindow struct {
	mutex sync.Mutex

	width  int32
	height int32
	zoom   int32

	commit bool

	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture

	buffer [][]Pixel

	keyboardCallback KeyboardCallback
}

func NewSDLWindow(title string, width, height int32) (*mainWindow, error) {
	zoom := int32(_defaultZoom)

	if err := sdl.InitSubSystem(_sdlSubSystemFlags); err != nil {
		return nil, err
	}

	windowWidth, windowHeight := windowSize(width, height, zoom)
	window, err := sdl.CreateWindow(title, _sdlWindowPos, _sdlWindowPos, windowWidth, windowHeight, uint32(_sdlWindowFlags))
	if err != nil {
		return nil, err
	}

	renderer, err := sdl.CreateRenderer(window, -1, _sdlRendererFlags)
	if err != nil {
		return nil, err
	}

	texture, err := renderer.CreateTexture(_sdlTextureFormat, _sdlTextureAccess, windowWidth, windowHeight)
	if err != nil {
		return nil, err
	}

	buffer := make([][]Pixel, height)
	for idx, _ := range buffer {
		buffer[idx] = make([]Pixel, width)
	}

	result := &mainWindow{
		width:  width,
		height: height,
		zoom:   zoom,

		window:   window,
		renderer: renderer,
		texture:  texture,

		buffer: buffer,

		keyboardCallback: nil,
	}

	return result, nil
}

func (m *mainWindow) MainLoop(ctx context.Context, cancel context.CancelFunc) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		default:
			if event := sdl.PollEvent(); event != nil {
				switch t := event.(type) {
				case *sdl.KeyboardEvent:
					switch t.Keysym.Sym {
					case sdl.K_ESCAPE:
						cancel()
						fallthrough

					default:
						m.callKeyboardCallback(t)
					}

				case *sdl.QuitEvent:
					cancel()
				}
			}
		}

		if err := m.Refresh(); err != nil {
			return err
		}

		sdl.Delay(1)
	}
}

func (m *mainWindow) callKeyboardCallback(event *sdl.KeyboardEvent) {
	if m.keyboardCallback == nil {
		return
	}

	abstractKey, ok := _sdlKeys[event.Keysym.Sym]
	if !ok {
		return
	}

	eventT := DownState
	if event.Type == sdl.KEYUP {
		eventT = UpState
	}

	m.keyboardCallback(abstractKey, eventT)
}

func (m *mainWindow) Close() error {
	err := multierr.Combine(
		m.texture.Destroy(),
		m.renderer.Destroy(),
		m.window.Destroy(),
	)

	if err != nil {
		return err
	}

	sdl.QuitSubSystem(_sdlSubSystemFlags)
	return nil
}

func (m *mainWindow) Commit() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.commit = true
	return nil
}

func (m *mainWindow) Refresh() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.commit {
		buffer, pitch := m.bufferAsBytes()

		err := multierr.Combine(
			m.renderer.Clear(),
			m.texture.Update(nil, buffer, pitch),
			m.renderer.Copy(m.texture, nil, nil),
		)
		if err != nil {
			return err
		}

		m.renderer.Present()
		m.commit = false
	}
	return nil
}

func (m *mainWindow) bufferAsBytes() ([]byte, int) {
	width, height := windowSize(m.width, m.height, m.zoom)
	screenB := make([]byte, 0, 4*width*height)

	for _, pixels := range m.buffer {
		lineB := make([]byte, 0, 4*width)

		for _, pixel := range pixels {
			pixelData := sdlData(pixel)

			for idx := int32(0); idx < m.zoom; idx++ {
				lineB = append(lineB, pixelData...)
			}
		}

		for idx := int32(0); idx < m.zoom; idx++ {
			screenB = append(screenB, lineB...)
		}
	}

	return screenB, rendererPitch(m.width, m.zoom)
}

func (m *mainWindow) DrawLine(line uint, pixels []Pixel) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(pixels) != int(m.width) {
		return errors.New("line is not the right width")
	}

	pixelsCopy := make([]Pixel, m.width)
	copy(pixelsCopy, pixels)
	m.buffer[line] = pixelsCopy

	return nil
}

func (m *mainWindow) SubscribeKeyboardStateChanges(callback KeyboardCallback) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.keyboardCallback = callback
}
