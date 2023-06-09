package vector3

import (
	"game-engine/math64"
	"math"
)

var Zero = Vector3{0, 0, 0}
var One = Vector3{1, 1, 1}

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

func (v Vector3) Max() float64 {
	var values = [3]float64{v.X, v.Y, v.Z}
	max := values[0]
	//var min float64 = values[0]
	for _, value := range values {
		if max < value {
			max = value
		}
		//if min > value {
		//	min = value
		//}
	}
	return max
}

func (v Vector3) DivScalar(scalar float64) Vector3 {
	return Vector3{
		v.X / scalar,
		v.Y / scalar,
		v.Z / scalar,
	}
}

func (v Vector3) Add(a Vector3) Vector3 {
	return Vector3{
		v.X + a.X,
		v.Y + a.Y,
		v.Z + a.Z,
	}
}

func (v Vector3) Sub(a Vector3) Vector3 {
	return Vector3{
		v.X - a.X,
		v.Y - a.Y,
		v.Z - a.Z,
	}
}

func (v Vector3) MulScalar(scalar float64) Vector3 {
	return Vector3{
		v.X * scalar,
		v.Y * scalar,
		v.Z * scalar,
	}
}

func (v Vector3) Normalize() Vector3 {
	length := v.Magnitude()
	if length > math64.Epsilon {
		return Vector3{
			X: v.X / length,
			Y: v.Y / length,
			Z: v.Z / length,
		}
	}
	return v
}

func (v Vector3) At(i int) float64 {
	if i == 0 {
		return v.X
	} else if i == 1 {
		return v.Y
	} else if i == 2 {
		return v.Z
	}
	panic("")
}

func (v Vector3) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector3) SqrMagnitude() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func Max(a Vector3, b Vector3) Vector3 {
	return Vector3{
		X: math.Max(a.X, b.X),
		Y: math.Max(a.Y, b.Y),
		Z: math.Max(a.Z, b.Z),
	}
}

func Min(a Vector3, b Vector3) Vector3 {
	return Vector3{
		X: math.Min(a.X, b.X),
		Y: math.Min(a.Y, b.Y),
		Z: math.Min(a.Z, b.Z),
	}
}

func Mul(a Vector3, b Vector3) Vector3 {
	return Vector3{
		a.X * b.X,
		a.Y * b.Y,
		a.Z * b.Z,
	}
}

func MulScalar(a Vector3, scalar float64) Vector3 {
	return Vector3{
		a.X * scalar,
		a.Y * scalar,
		a.Z * scalar,
	}
}

func Sum(a Vector3, b Vector3) Vector3 {
	return Vector3{
		a.X + b.X,
		a.Y + b.Y,
		a.Z + b.Z,
	}
}

func SumScalar(a Vector3, scalar float64) Vector3 {
	return Vector3{
		a.X + scalar,
		a.Y + scalar,
		a.Z + scalar,
	}
}

func Sub(a Vector3, b Vector3) Vector3 {
	return Vector3{
		a.X - b.X,
		a.Y - b.Y,
		a.Z - b.Z,
	}
}

func Div(a Vector3, b Vector3) Vector3 {
	return Vector3{
		a.X / b.X,
		a.Y / b.Y,
		a.Z / b.Z,
	}
}

func DivScalar(scalar float64, a Vector3) Vector3 {
	return Vector3{
		scalar / a.X,
		scalar / a.Y,
		scalar / a.Z,
	}
}
