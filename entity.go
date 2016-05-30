package goldengine

import (
	"sync/atomic"
	"time"
)

//EntityPrefab : Information Required to Create an Entity from JSON Prefab
type EntityPrefab struct {
	Name        string
	Components  []ComponentPrefab
	Transformer TransformerPrefab
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

//RemoveChild : Removes Child Enttiy
func (e *Entity) RemoveChild(child *Entity) {
	delete(e.children, child.id)
}

//AddComponent : Associates a component with a given entity
func (e *Entity) AddComponent(comp Component) {
	if comp != nil {
		e.components = append(e.components, comp)
		comp.SetEntity(e)
	}
}

var entityCounter uint32 = 1

//NewEntity : Create a new Entity
func NewEntity() *Entity {
	id := entityCounter
	atomic.AddUint32(&entityCounter, 1)
	return &Entity{
		id:          id,
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
	e.Transfrom.Rotate(1.0)

	e.components = make([]Component, 0)
	for _, p := range prefab.Components {

		comp := ComponentFromComponentPrefab(p)
		e.AddComponent(comp)

	}
	return e
}
