package main

import (
	"fmt"
	"time"

	GE "github.com/Dacode45/goldengine"
	sf "github.com/manyminds/gosfml"
)

func init() {
	GE.ComponentRegister.Register("paddle", NewPaddleComponent)
}

//PaddleComponent : Moves the Paddle Up and Down when Up and down keys are pressed
type PaddleComponent struct {
	direction  float32
	Speed      float32
	keyHandler *GE.KeyboardHandler
	GE.BaseComponent
}

//NewPaddleComponent : ComponentGenerator for Paddle
func NewPaddleComponent(args map[string]interface{}) GE.Component {
	comp := PaddleComponent{
		keyHandler: GE.GenInputHandler(),
	}
	comp.keyHandler.RegisterKeyPressedCommand(sf.KeyUp, comp.MoveUp)
	comp.keyHandler.RegisterKeyPressedCommand(sf.KeyDown, comp.MoveDown)
	comp.keyHandler.RegisterKeyReleasedCommand(sf.KeyUp, comp.StopMovement)
	comp.keyHandler.RegisterKeyReleasedCommand(sf.KeyDown, comp.StopMovement)
	if arg, ok := args["Speed"]; ok {
		if speed, ok := GE.ArgAsFloat32(arg); ok {
			comp.Speed = speed
		} else {
			comp.Speed = 400.0
		}
	}
	return &comp
}

//MoveUp : Moves Paddle Up
func (comp *PaddleComponent) MoveUp() {
	comp.direction = -1
}

//StopMovement : Paddle Won't move anymore
func (comp *PaddleComponent) StopMovement() {
	fmt.Println("Stopped Movement")
	comp.direction = 0
}

//MoveDown : Moves the Paddle Down
func (comp *PaddleComponent) MoveDown() {
	comp.direction = 1
}

//Awake : Component now listening for keyboardinput
func (comp *PaddleComponent) Awake() {
	fmt.Println("Paddle Awake")
	comp.GetEntity().KeyboardSet.AddHandler(comp.keyHandler)
}

//Update : Moves the Paddle
func (comp *PaddleComponent) Update(dur time.Duration) {
	transform := comp.GetEntity().Transfrom
	distance := sf.Vector2f{
		X: 0,
		Y: comp.direction * comp.Speed * float32(dur.Seconds()),
	}
	transform.Move(distance)
}

//Sleep : Component stops listening for input
func (comp *PaddleComponent) Sleep() {
	comp.GetEntity().KeyboardSet.RemoveHandler(comp.keyHandler)
}
