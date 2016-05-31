package main

import (
	"fmt"

	GE "github.com/Dacode45/goldengine"
	sf "github.com/manyminds/gosfml"
)

func main() {
	app := GE.NewGame(GE.GameConfig{Name: "Test", Debug: true}, GE.WindowConfig{
		Width:      800,
		Height:     600,
		ClearColor: sf.Color{R: 50, G: 200, B: 50, A: 0},
	}, GE.PhysicsEngineConfig{Debug: true},
	)
	app.ProcessArguments()
	app.Init()
	app.ChangeScene("main")
	scene := app.GetCurrentScene()
	scene.Start = func() {
		fmt.Println("In Start")
		paddle, _ := scene.GetEntityByName("leftPaddle")
		app.GetWindow().GetInputCollection().InstallKeyboardSet(paddle.KeyboardSet)
	}
	app.Run()
}
