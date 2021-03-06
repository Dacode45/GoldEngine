package goldengine

import (
	"fmt"
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

//Window : Wrapper around a SFML RenderWindow Handles input and drawing scenes
type Window struct {
	Ticker          *time.Ticker
	ClearColor      sf.Color
	renderWindow    *sf.RenderWindow
	scene           *Scene
	inputCollection *InputCollection

	BasicMailBox
}

const (
	//DefaultGameWidth : Width in Pixels of Game Window
	DefaultGameWidth = 800
	//DefaultGameHeight : Height in Pixels of Game Window
	DefaultGameHeight = 600
	//DefaultTitle : Title of Game Window
	DefaultTitle = "Gold Engine"
)

func newWindow(config WindowConfig) *Window {
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
	NewScreenWidth(gameWidth)
	w := &Window{
		renderWindow:    sf.NewRenderWindow(sf.VideoMode{Width: gameWidth, Height: gameHeight, BitsPerPixel: 32}, gameTitle, sf.StyleDefault, sf.DefaultContextSettings()),
		Ticker:          time.NewTicker(time.Second / 60),
		ClearColor:      config.ClearColor,
		inputCollection: GenInputCollection(),
	}
	return w
}

//GetInputCollection : Retruns the InputCollection for this scene
func (w *Window) GetInputCollection() *InputCollection {
	return w.inputCollection
}

//ChangeScene : Changes current scene
func (w *Window) ChangeScene(s *Scene) {
	w.scene = s
}

//Init : Subscribes to messages and other Stuff
func (w *Window) Init() {
	//SetSubscriptions
	w.GetOffice().Add(w.inputCollection)
	w.GetOffice().Subscribe(w.GetAddress(), w.inputCollection.GetAddress(), KeyPressedMSG)
	w.GetOffice().Subscribe(w.GetAddress(), w.inputCollection.GetAddress(), KeyReleasedMSG)
}

//RecieveMessage : Messages This Window should Recieve
func (w *Window) RecieveMessage(msg Message) {
	switch msg.Message {
	case SceneChangedMSG:
		w.ChangeScene(msg.Content.(*Scene))
	}
}

//Run : Plays the window
func (w *Window) Run() {
	if w.scene == nil {
		panic(fmt.Errorf("No Scene"))
	}
	for w.renderWindow.IsOpen() {
		select {
		case <-w.Ticker.C:
			for event := w.renderWindow.PollEvent(); event != nil; event = w.renderWindow.PollEvent() {
				switch ev := event.(type) {
				case sf.EventKeyPressed:
					var code = ev.Code
					msg := Message{KeyPressedMSG, &code}
					w.PostMessage(msg)
				case sf.EventKeyReleased:
					var code = ev.Code
					msg := Message{KeyReleasedMSG, &code}
					w.PostMessage(msg)
				case sf.EventClosed:
					w.renderWindow.Close()
				case sf.EventResized:
					NewScreenWidth(ev.Width)
				}
			}

			w.renderWindow.Clear(w.ClearColor)
			w.renderWindow.Draw(w.scene, sf.DefaultRenderStates())

			w.renderWindow.Display()
		}
	}
}
