package geo

type Sphere struct {
	Shape
	Radius float64
}

func (s *Sphere) Type() ShapeType {
	return ShapeSphere
}
