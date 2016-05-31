package goldengine

import (
	"time"

	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
)

import sf "github.com/manyminds/gosfml"

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
}

//ChangeScene : Changes the Current Scene
func (engine *PhysicsEngine) ChangeScene(s *Scene) {
	//TODO : Remove all debug components
	engine.scene = s
}

//AddEntity : Add Entity to phsics space
func (engine *PhysicsEngine) AddEntity(e *Entity) {
	if e != nil && e.Collider != nil {
		engine.space.AddBody(e.Collider)
	}
	if engine.debug {

	}
}

func (engine *PhysicsEngine) makeDebugEntity(e *Entity) {
	if engine.scene == nil {
		return
	}
	debugRoot := NewEntity()
	for _, shape := range e.Collider.Shapes {

		switch shapeClass := shape.ShapeClass.(type) {
		case *chipmunk.CircleShape:
			circle, err := sf.NewCircleShape()
			if err != nil {
				continue
			}
			circle.SetRadius(float32(shapeClass.Radius))
			circle.SetFillColor(sf.ColorYellow())
			debug := NewEntity()
			debug.Transfrom = circle
			debugRoot.AddChild(debug)
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
