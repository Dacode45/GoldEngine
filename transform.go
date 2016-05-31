package goldengine

import sf "github.com/manyminds/gosfml"

//Transformer : Wrapper arround sfml Transformer
type Transformer interface {
	sf.Drawer
	sf.Transformer
}

//TransformerPrefab : Info to create a Transformer from JSON Prefab
type TransformerPrefab struct {
	Kind      string
	Arguments map[string]interface{}
}

//Sprite : Wrapper around sfml Object
type Sprite sf.Sprite

//CircleShape : Wrapper around sfml Object
type CircleShape sf.CircleShape

//ConvexShape : Wrapper around sfml Object
type ConvexShape sf.ConvexShape

//RectangleShape : Wrapper around sfml Object
type RectangleShape sf.RectangleShape

//Text : Wrapper around sfml Object
type Text sf.Text

//Shape : Group of Functions all Shapes have
type Shape interface {
	GetOrigin() sf.Vector2f
	SetOrigin(sf.Vector2f)

	GetOutlineThickness() float32
	SetOutlineThickness(float32)

	GetOutlineColor() sf.Color
	SetOutlineColor(sf.Color)

	GetFillColor() sf.Color
	SetFillColor(sf.Color)
}

const (
	//SpriteName : Name of Transformer
	SpriteName = "Sprite"
	//CircleShapeName : Name of Transformer
	CircleShapeName = "CircleShape"
	//ConvexShapeName : Name of Transformer
	ConvexShapeName = "ConvexShape"
	//RectangleShapeName : Name of Transformer
	RectangleShapeName = "RectangleShape"
	//TextName : Name of Transformer
	TextName = "Text"
)

//TranformerGenerators : map to create Generators
var TranformerGenerators = map[string]func(args map[string]interface{}) (Transformer, error){
	SpriteName:         SpriteFromArguments,
	CircleShapeName:    CircleShapeFromArguments,
	ConvexShapeName:    ConvexShapeFromArguments,
	RectangleShapeName: RectangleShapeFromArguments,
	TextName:           TextFromArguments,
}

//TransformerFromTranformerPrefab : Returns Transformer from Transform Prefab
func TransformerFromTranformerPrefab(t TransformerPrefab) (Transformer, error) {
	var transformer Transformer
	generator, ok := TranformerGenerators[t.Kind]
	if !ok {
		return nil, nil
	}
	transformer, err := generator(t.Arguments)
	return transformer, err
}

/*SpriteFromArguments : Generates sprite from Arguments field of Prefab
TODO : Get Texture*/
func SpriteFromArguments(args map[string]interface{}) (Transformer, error) {
	return sf.NewSprite(nil)
}

//CircleShapeFromArguments : Generates CircleShape from Arguments field of Prefab
func CircleShapeFromArguments(args map[string]interface{}) (Transformer, error) {
	shape, err := sf.NewCircleShape()
	if err != nil {
		return nil, err
	}
	ApplyArgsToShape(shape, args)
	if arg, ok := args["Radius"]; ok {
		radius, ok := ArgAsFloat32(arg)
		if ok {
			shape.SetRadius(radius)
		}
	}
	return shape, err
}

//ConvexShapeFromArguments : Generates ConvexShape from Arguments field of Prefab
func ConvexShapeFromArguments(args map[string]interface{}) (Transformer, error) {
	return sf.NewConvexShape()
}

//RectangleShapeFromArguments : Generates RectangleShape from Arguments field of Prefab
func RectangleShapeFromArguments(args map[string]interface{}) (Transformer, error) {
	shape, err := sf.NewRectangleShape()
	if err != nil {
		return nil, err
	}
	ApplyArgsToShape(shape, args)
	if arg, ok := args["Size"]; ok {
		size, ok := ArgAsVector2f(arg)
		if ok {
			shape.SetSize(size)
		}
	}
	return shape, err
}

//ApplyArgsToShape : Sets Properties like OutlineThickness
func ApplyArgsToShape(shape Shape, args map[string]interface{}) {
	if arg, ok := args["OutlineThickness"]; ok {
		thickness, ok := ArgAsFloat32(arg)
		if ok {
			shape.SetOutlineThickness(thickness)
		}
	}

	if arg, ok := args["OutlineColor"]; ok {
		color, ok := ArgAsColor(arg)
		if ok {
			shape.SetOutlineColor(color)
		}
	}

	if arg, ok := args["Origin"]; ok {
		origin, ok := ArgAsVector2f(arg)
		if ok {
			shape.SetOrigin(origin)
		}
	}

	if arg, ok := args["FillColor"]; ok {
		color, ok := ArgAsColor(arg)
		if ok {
			shape.SetFillColor(color)
		}
	}

}

/*TextFromArguments : Generates Text from Arguments field of Prefab
TODO : Get fonts*/
func TextFromArguments(args map[string]interface{}) (Transformer, error) {
	return sf.NewText(nil)
}

//ArgAsFloat32 Converts an interface from a JSON Parser to a float32
func ArgAsFloat32(arg interface{}) (float32, bool) {
	value, ok := arg.(float64)
	if ok {
		return float32(value), ok
	}
	return 0.0, ok
}

//ArgAsVector2f Converts an interface from a JSON Parser to a Vector2f
func ArgAsVector2f(arg interface{}) (sf.Vector2f, bool) {
	value, ok := arg.(map[string]interface{})
	if ok {
		vec, ok := Vector{X: float32(value["X"].(float64)), Y: float32(value["Y"].(float64))}, ok
		return vec.ToSFML(), ok
	}
	return sf.Vector2f{}, ok
}

//ArgAsVector Converts an interface from a JSON Parser to a Vector2f
func ArgAsVector(arg interface{}) (Vector, bool) {
	value, ok := arg.(map[string]interface{})
	if ok {
		vec, ok := Vector{X: float32(value["X"].(float64)), Y: float32(value["Y"].(float64))}, ok
		return vec, ok
	}
	return Vector{}, ok
}

//ArgAsColor Converts an interface from a JSON Parser to a Color
func ArgAsColor(arg interface{}) (sf.Color, bool) {
	value, ok := arg.(map[string]interface{})
	if ok {
		return sf.Color{
			R: uint8(value["R"].(float64)),
			G: uint8(value["G"].(float64)),
			B: uint8(value["B"].(float64)),
			A: uint8(value["A"].(float64)),
		}, ok
	}
	return sf.Color{}, ok
}
