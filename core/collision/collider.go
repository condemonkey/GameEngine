package collision

type ShapeType int

const (
	ShapeSphere  ShapeType = 1
	ShapeCapsule           = 2
	ShapeBox               = 3
)

type Collider interface {
	Type() ShapeType
	FatAABB() *AABB
	AABB() *AABB
}
