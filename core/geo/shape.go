package geo

type ShapeType int

const (
	ShapeSphere  ShapeType = 1
	ShapeCapsule           = 2
	ShapeBox               = 3
	ShapePlane             = 4
)

type Shape interface {
	Type() ShapeType
}
