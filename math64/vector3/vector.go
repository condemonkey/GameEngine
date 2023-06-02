package vector3

import "math"

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

func (v Vector3) Scale(a Vector3) Vector3 {
	return Vector3{
		v.X - a.X,
		v.Y - a.Y,
		v.Z - a.Z,
	}
}

func (v Vector3) MulScala(a float64) Vector3 {
	return Vector3{
		v.X * a,
		v.Y * a,
		v.Z * a,
	}
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

func Scale(a Vector3, b Vector3) Vector3 {
	return Vector3{
		a.X * b.X,
		a.Y * b.Y,
		a.Z * b.Z,
	}
}

func Scalef(a Vector3, s float64) Vector3 {
	return Vector3{
		a.X * s,
		a.Y * s,
		a.Z * s,
	}
}

func Sum(a Vector3, b Vector3) Vector3 {
	return Vector3{
		a.X + b.X,
		a.Y + b.Y,
		a.Z + b.Z,
	}
}

func Sumf(a Vector3, f float64) Vector3 {
	return Vector3{
		a.X + f,
		a.Y + f,
		a.Z + f,
	}
}

func Sub(a Vector3, b Vector3) Vector3 {
	return Vector3{
		a.X - b.X,
		a.Y - b.Y,
		a.Z - b.Z,
	}
}
