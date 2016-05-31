package goldengine

import (
	"errors"

	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
)

//ColliderPrefab : Generates a Collider (chipmunk body) from configuration
type ColliderPrefab struct {
	Kind      string
	Arguments map[string]interface{}
}

const (
	//CircleColliderName : Name of Circle Collider
	CircleColliderName = "CircleShape"
)

//ColliderGenerators : map to generators for bodies
var ColliderGenerators = map[string]func(args map[string]interface{}) (*chipmunk.Body, error){
	CircleColliderName: CircleColliderFromColliderPrefab,
}

//ColliderFromColliderPrefab : Returns a Collider from ColliderPrefab
func ColliderFromColliderPrefab(c ColliderPrefab) (*chipmunk.Body, error) {
	var collider *chipmunk.Body
	generator, ok := ColliderGenerators[c.Kind]
	if !ok {
		return nil, errors.New("No collider with that name")
	}
	collider, err := generator(c.Arguments)
	return collider, err
}

//CircleColliderFromColliderPrefab : Creates a body shape
func CircleColliderFromColliderPrefab(args map[string]interface{}) (*chipmunk.Body, error) {
	var shape *chipmunk.Shape
	if arg, ok := args["Radius"]; ok {
		radius, ok := ArgAsFloat32(arg)
		if !ok {
			radius = 1
		}
		shape = chipmunk.NewCircle(vect.Vector_Zero, radius)
	} else {
		shape = chipmunk.NewCircle(vect.Vector_Zero, 1)
	}
	body := chipmunk.NewBody(1, shape.Moment(1))
	ApplyChipmunkShapeProperties(shape, args)
	ApplyChipmunkBodyProperties(shape, body, args)
	body.AddShape(shape)
	return body, nil
}

//ApplyChipmunkShapeProperties : Sets common chipmunk shape properties
func ApplyChipmunkShapeProperties(shape *chipmunk.Shape, args map[string]interface{}) {
	if arg, ok := args["Elasticity"]; ok {
		elasticity, ok := ArgAsChipmunkFloat(arg)
		if ok {
			shape.SetElasticity(elasticity)
		}
	}
}

//ApplyChipmunkBodyProperties : Sets common chipmunk body properties
func ApplyChipmunkBodyProperties(shape *chipmunk.Shape, body *chipmunk.Body, args map[string]interface{}) {
	if arg, ok := args["Mass"]; ok {
		mass, ok := ArgAsChipmunkFloat(arg)
		if ok {
			body.SetMass(mass)
		} else {
			body.SetMass(vect.Float(1))
		}
	}

	if arg, ok := args["Moment"]; ok {
		moment, ok := ArgAsFloat32(arg)
		if ok {
			body.SetMoment(shape.Moment(moment))
		} else {
			body.SetMoment(shape.Moment(float32(body.Mass())))
		}
	}

	if arg, ok := args["Position"]; ok {
		pos, ok := ArgAsChipmunkVector(arg)
		if ok {
			body.SetPosition(pos)
		}
	}
}

//ArgAsChipmunkVector : Converts an interface from a JSON Parser to a chipmunk
func ArgAsChipmunkVector(arg interface{}) (vect.Vect, bool) {
	value, ok := arg.(map[string]interface{})
	if ok {
		vec, ok := Vector{X: float32(value["X"].(float64)), Y: float32(value["Y"].(float64))}, ok
		return vec.ToChipmunk(), ok
	}
	return vect.Vect{}, ok
}

//ArgAsChipmunkFloat : Converts an interface from a JSON Parser to a chipmunk Float
func ArgAsChipmunkFloat(arg interface{}) (vect.Float, bool) {
	value, ok := arg.(float64)
	if ok {
		return vect.Float(value), ok
	}
	return 0.0, ok
}
