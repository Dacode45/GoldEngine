package goldengine

import "sync/atomic"

//EntityPrefab : Information Required to Create an Entity from JSON Prefab
type EntityPrefab struct {
	Name        string
	Components  []ComponentPrefab
	Transformer TransformerPrefab
}

//Entity : Entity
type Entity struct {
	Name       string
	parent     *Entity
	children   map[uint32]*Entity
	id         uint32
	components []Component
	Transfrom  Transformer
	tags       []string
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
		id: id,
	}
}

//EntityFromEntityPrefab : Returns Entity from Entity Prefab
func EntityFromEntityPrefab(prefab EntityPrefab) *Entity {
	e := NewEntity()
	e.Name = prefab.Name
	e.Transfrom, _ = TransformerFromTranformerPrefab(prefab.Transformer)
	e.components = make([]Component, len(prefab.Components))
	for _, p := range prefab.Components {

		comp := ComponentFromComponentPrefab(p)
		e.AddComponent(comp)

	}
	return e
}
