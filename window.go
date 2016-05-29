package goldengine

import (
	"time"

	sf "github.com/manyminds/gosfml"
)

//WindowConfig : Creates A Window given the configuration
type WindowConfig struct {
	Width      uint
	Height     uint
	ClearColor sf.Color
	Title      string
}

type window struct {
	Ticker       *time.Ticker
	ClearColor   sf.Color
	renderWindow *sf.RenderWindow
	scene        *Scene
}

const (
	//DefaultGameWidth : Width in Pixels of Game Window
	DefaultGameWidth = 800
	//DefaultGameHeight : Height in Pixels of Game Window
	DefaultGameHeight = 600
	//DefaultTitle : Title of Game Window
	DefaultTitle = "Gold Engine"
)

func newWindow(config WindowConfig) *window {
	gameWidth := config.Width
	if gameWidth <= 0 {
		gameWidth = DefaultGameWidth
	}

	gameHeight := config.Height
	if gameHeight <= 0 {
		gameHeight = DefaultGameHeight
	}

	gameTitle := config.Title
	if gameTitle == "" {
		gameTitle = DefaultTitle
	}

	return &window{
		renderWindow: sf.NewRenderWindow(sf.VideoMode{Width: gameWidth, Height: gameHeight, BitsPerPixel: 32}, gameTitle, sf.StyleDefault, sf.DefaultContextSettings()),
		Ticker:       time.NewTicker(time.Second / 60),
		ClearColor:   config.ClearColor,
	}
}

func (w *window) Run() {
	for w.renderWindow.IsOpen() {
		select {
		case <-w.Ticker.C:
			w.renderWindow.Clear(w.ClearColor)
			if w.scene != nil {
				w.renderWindow.Draw(w.scene, sf.DefaultRenderStates())
			}
			w.renderWindow.Display()
		}
	}
}
