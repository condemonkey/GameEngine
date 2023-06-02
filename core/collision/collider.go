package collision

type ColliderType int

const (
	ShapeSphere  ColliderType = 1
	ShapeCapsule              = 2
	ShapeBox                  = 3
)

type Collider interface {
	Type() ColliderType
	FatAABB() *AABB
	AABB() *AABB
	Raycast(ray *Ray, maxDistance float64, hit *RaycastHit) bool
}
