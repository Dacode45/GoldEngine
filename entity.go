package goldengine

import (
	"strconv"
	"sync/atomic"
	"time"

	sf "github.com/manyminds/gosfml"
)

import "github.com/vova616/chipmunk"

//EntityPrefab : Information Required to Create an Entity from JSON Prefab
type EntityPrefab struct {
	Name        string
	Components  []ComponentPrefab
	Transformer TransformerPrefab
	Collider    ColliderPrefab
}

//Entity : Entity
type Entity struct {
	Name        string
	id          uint32
	tags        []string
	KeyboardSet *KeyboardSet
	parent      *Entity
	children    map[uint32]*Entity
	components  []Component
	Transfrom   Transformer
	Collider    *chipmunk.Body
	scene       *Scene
	started     bool
	awake       bool
}

//Start : Called before gameloop
func (e *Entity) Start() {
	e.started = true
	for _, c := range e.components {
		c.Start()
	}
}

//Awake : Called After Start and after sleep
func (e *Entity) Awake() {
	e.awake = true
	for _, c := range e.components {
		c.Awake()
	}
}

//Update : Called during the gameloop
func (e *Entity) Update(dur time.Duration) {
	for _, c := range e.components {
		c.Update(dur)
	}
}

//Sleep : Shouldn't Update Entity Anymore
func (e *Entity) Sleep() {
	e.awake = false
	for _, c := range e.components {
		c.Sleep()
	}
}

//Stop : Shouldn't Update Entity Anymore
func (e *Entity) Stop() {
	e.started = false
	for _, c := range e.components {
		c.Stop()
	}
}

//AddChild : Adds A child to the entity
func (e *Entity) AddChild(child *Entity) {
	child.parent = e
	e.children[child.id] = child
}

//RemoveChild : Removes Child Enttiy
func (e *Entity) RemoveChild(child *Entity) {
	child.parent = nil
	delete(e.children, child.id)
}

//AddComponent : Associates a component with a given entity
func (e *Entity) AddComponent(comp Component) {
	if comp != nil {
		e.components = append(e.components, comp)
		comp.SetEntity(e)
	}
}

//RecalculateScale : Changes the size of sfml and chipmunk objects
func (e *Entity) RecalculateScale() {

}

var entityCounter uint32 = 1

//NewEntity : Create a new Entity
func NewEntity() *Entity {
	id := entityCounter
	atomic.AddUint32(&entityCounter, 1)
	return &Entity{
		id:          id,
		Name:        strconv.Itoa(int(id)),
		KeyboardSet: GenKeyboardSet(),
	}
}

//EntityFromEntityPrefab : Returns Entity from Entity Prefab
func EntityFromEntityPrefab(prefab EntityPrefab) *Entity {
	e := NewEntity()
	e.Name = prefab.Name
	var err error
	e.Transfrom, err = TransformerFromTranformerPrefab(prefab.Transformer)
	if err != nil {
		panic(err)
	}
	if prefab.Collider.Kind != "" {
		prefab.Collider = ColliderPrefabFromTransformPrefab(prefab.Collider, prefab.Transformer, e.Transfrom)
		e.Collider, err = ColliderFromColliderPrefab(prefab.Collider)
		if err != nil {
			panic(err)
		}
	}
	e.components = make([]Component, 0)
	for _, p := range prefab.Components {

		comp := ComponentFromComponentPrefab(p)
		e.AddComponent(comp)

	}
	return e
}

//ColliderPrefabFromTransformPrefab : Copies arguments from Transforms
func ColliderPrefabFromTransformPrefab(collider ColliderPrefab, transform TransformerPrefab, transformer Transformer) ColliderPrefab {
	if transform.Kind == collider.Kind {
		if _, has := collider.Arguments["Radius"]; !has {
			if _, ok := transform.Arguments["Radius"]; ok {
				collider.Arguments["Radius"] = transform.Arguments["Radius"]
			} else {
				shape, ok := transformer.(*sf.CircleShape)
				if !ok {
					panic("Not a circle shape")
				}
				collider.Arguments["Radius"] = float64(shape.GetRadius() / scale)
			}
		}
	}
	if _, has := collider.Arguments["Position"]; !has {
		if _, ok := transform.Arguments["Position"]; ok {
			collider.Arguments["Position"] = transform.Arguments["Position"]
		} else {
			vec := Vector2fToVector(transformer.GetPosition())
			collider.Arguments["Position"] = map[string]float64{
				"X": float64(vec.X), "Y": float64(vec.Y),
			}
		}
	}
	if _, has := collider.Arguments["Size"]; !has {
		if _, ok := transform.Arguments["Size"]; ok {
			collider.Arguments["Size"] = transform.Arguments["Size"]
		} else {
			shape := transformer.(*sf.RectangleShape)
			vec := Vector2fToVector(shape.GetSize())
			collider.Arguments["Size"] = map[string]float64{
				"X": float64(vec.X), "Y": float64(vec.Y),
			}
		}
	}
	return collider
}
