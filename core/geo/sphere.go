package geo

type Sphere struct {
	Radius float64
}

func NewSphere(radius float64) Sphere {
	return Sphere{
		Radius: radius,
	}
}

func (s Sphere) Type() ShapeType {
	return ShapeSphere
}
