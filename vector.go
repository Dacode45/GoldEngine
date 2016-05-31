package goldengine

import (
	"fmt"

	sf "github.com/manyminds/gosfml"
	"github.com/vova616/chipmunk/vect"
)

//Vector : A 1x1 Vector reprecents a square whose width is 1 100th of the screen
type Vector struct {
	X, Y float32
}

//ZeroVector : Vector initialized to zero
var ZeroVector = Vector{X: 0, Y: 0}

//Assumed an 800x600 resolution
var oldscale float32 = 8
var scale float32 = 8

//NewScreenWidth : Change The screenWidth
func NewScreenWidth(width uint) {
	oldscale = scale
	scale = float32(width) / 100
	fmt.Printf("oldscale : %v, scale : %v\n", oldscale, scale)
}

//ToSFML : Converts a Vector to an sfml vector
func (v *Vector) ToSFML() sf.Vector2f {
	return sf.Vector2f{
		X: v.X * scale,
		Y: v.Y * scale,
	}
}

//ToChipmunk : Converts a Vector to an sfml vector
func (v *Vector) ToChipmunk() vect.Vect {
	return vect.Vect{
		X: vect.Float(v.X * scale),
		Y: vect.Float(v.Y * scale),
	}
}

//Vector2fToVector : Converts a Vector2f to a Vector
func Vector2fToVector(vec sf.Vector2f) Vector {
	return Vector{
		X: vec.X / scale,
		Y: vec.Y / scale,
	}
}

//Vector2uToVector : Converts a Vector2u to a Vector
func Vector2uToVector(vec sf.Vector2u) Vector {
	return Vector{
		X: float32(vec.X) / scale,
		Y: float32(vec.Y) / scale,
	}
}

//ChipmunkVectorToVector : New Vector
func ChipmunkVectorToVector(vec vect.Vect) Vector {
	return Vector{
		X: float32(vec.X) / scale,
		Y: float32(vec.Y) / scale,
	}
}

//ChipmunkFloatToVector : New float
func ChipmunkFloatToVector(x, y vect.Float) Vector {
	return Vector{
		X: float32(x) / scale,
		Y: float32(y) / scale,
	}
}
