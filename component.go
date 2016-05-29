package goldengine

import (
	"errors"
	"time"
)

// Component : Grants an entity an ability
type Component interface {
	Awake()
	Start()
	Update(time.Duration)
	Sleep()
	SetEntity(e *Entity)
	GetEntity() *Entity
}

// ComponentPrefab : Generates a Component with the Given name from the Arguments
type ComponentPrefab struct {
	Name      string
	Arguments map[string]interface{}
}

// ComponentFromComponentPrefab : Returns a Component from a ComponentPrefab
func ComponentFromComponentPrefab(c ComponentPrefab) Component {
	generator, ok := ComponentRegister.Get(c.Name)
	if !ok {
		return nil
	}
	return generator(c.Arguments)
}

// ComponentGenerator : Function that takes the minimum arguments required to create a component
type ComponentGenerator func(args map[string]interface{}) Component

// ComponentRegister : Registers a Generator With a ComponentName
type componentRegister struct {
	register map[string]ComponentGenerator
}

// ComponentRegister : map of generators for a componentName
var ComponentRegister = componentRegister{
	register: make(map[string]ComponentGenerator),
}

// ErrComponentAlreadyRegistered : Cannot register a ComponentGenerator because the name is already registered
var ErrComponentAlreadyRegistered = errors.New("A Component With That Name Was Already Registered")

// Register : Adds a ComponentGenerator for a component name
func (c *componentRegister) Register(name string, generator ComponentGenerator) error {
	_, ok := c.register[name]
	if ok {
		return ErrComponentAlreadyRegistered
	}
	c.register[name] = generator
	return nil
}

func (c *componentRegister) UnRegister(name string) {
	delete(c.register, name)
}

func (c *componentRegister) Get(name string) (ComponentGenerator, bool) {
	gen, ok := c.register[name]
	return gen, ok
}
