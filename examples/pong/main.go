package main

import GE "github.com/Dacode45/goldengine"
import sf "github.com/manyminds/gosfml"

func main() {
	app := GE.NewGame(GE.GameConfig{Name: "Test"}, GE.WindowConfig{
		Width:      800,
		Height:     600,
		ClearColor: sf.Color{R: 50, G: 200, B: 50, A: 0},
	})
	app.ProcessArguments()
	app.Init()
	app.ChangeScene("main")
	app.Run()
}
