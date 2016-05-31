package goldengine

import (
	"fmt"
	"time"

	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"

	sf "github.com/manyminds/gosfml"
)

//PhysicsEngineConfig : Configuration for the Physics Engine
type PhysicsEngineConfig struct {
	Debug   bool
	Gravity vect.Vect
}

//PhysicsEngine : Wrapper around chipmunk engine
type PhysicsEngine struct {
	debug         bool
	debugEntities []*Entity
	space         *chipmunk.Space
	scene         *Scene
	BasicMailBox
}

//RecieveMessage : RecieveMessage function for Physics Engine
func (engine *PhysicsEngine) RecieveMessage(msg Message) {
	switch msg.Message {
	case SceneChangedMSG:
		engine.ChangeScene(msg.Content.(*Scene))
	case SceneAddedEntityMSG:
		engine.AddEntity(msg.Content.(*Entity))
	}
}

//ChangeScene : Changes the Current Scene
func (engine *PhysicsEngine) ChangeScene(s *Scene) {
	//TODO : Remove all debug components

	office := engine.GetOffice()
	office.Subscribe(s.GetAddress(), engine.GetAddress(), SceneAddedEntityMSG)
	engine.scene = s
	engine.AddEntities(s.GetEntities())
}

//AddEntity : Add Entity to phsics space
func (engine *PhysicsEngine) AddEntity(e *Entity) {
	if e != nil && e.Collider != nil {
		engine.space.AddBody(e.Collider)
	}
	if engine.debug {
		engine.makeDebugEntity(e)
	}
}

//AddEntities : Convienece function to add a bunch of entities
func (engine *PhysicsEngine) AddEntities(list []*Entity) {
	for _, e := range list {
		engine.AddEntity(e)
	}
}

func (engine *PhysicsEngine) makeDebugEntity(e *Entity) {
	if engine.scene == nil || e.Collider == nil {
		fmt.Println(e.Collider)
		return
	}
	fmt.Println("Adding Entity", e.Name)
	debugRoot := NewEntity()
	for _, shape := range e.Collider.Shapes {

		switch shapeClass := shape.ShapeClass.(type) {
		case *chipmunk.CircleShape:
			circle, err := sf.NewCircleShape()
			if err != nil {
				continue
			}
			fmt.Println(shapeClass.Radius)
			circle.SetRadius(float32(shapeClass.Radius))
			// pos := e.Collider.Position()
			// vec := ChipmunkVectorToVector(pos)
			// circle.SetPosition(vec.ToSFML())
			circle.SetPosition(sf.Vector2f{X: 20, Y: 20})
			circle.SetFillColor(sf.ColorYellow())
			debug := NewEntity()
			debug.Transfrom = circle
			debugRoot.AddChild(debug)
			engine.scene.AddEntity(debug)
			engine.scene.SetZIndex(debug.Name, 1)
		case *chipmunk.BoxShape:
			rectangle, err := sf.NewRectangleShape()
			if err != nil {
				continue
			}
			box := shape.GetAsBox()
			size := ChipmunkFloatToVector(box.Width, box.Height)
			rectangle.SetSize(size.ToSFML())
			rectangle.SetFillColor(sf.ColorYellow())
			debug := NewEntity()
			debug.parent = debugRoot
			debug.Transfrom = rectangle
			debugRoot.AddChild(debug)
		}
	}
	engine.scene.AddEntity(debugRoot)
}

//SetDebug : Sets Debug mode of Physics Engine
func (engine *PhysicsEngine) SetDebug(debug bool) {
	engine.debug = debug
}

//GetSpace : Get chipmunk space
func (engine *PhysicsEngine) GetSpace() *chipmunk.Space {
	return engine.space
}

func (engine *PhysicsEngine) step(dur time.Duration) {
	engine.space.Step(vect.Float(dur.Seconds()))
}

//PhysicsEngineFromConfig : Generates a PhysicsEngine from Config
func newPhysicsEngine(config PhysicsEngineConfig) *PhysicsEngine {
	engine := &PhysicsEngine{
		debug: config.Debug,
		space: chipmunk.NewSpace(),
	}
	engine.space.Gravity = config.Gravity
	return engine
}
