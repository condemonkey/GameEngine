package collision

import (
	"game-engine/core/geo"
	"game-engine/math64"
	"game-engine/math64/vector3"
)

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
	IntersectRay(ray *Ray, maxDistance float64, hit *RaycastHit) bool
	IntersectShape(other Collider) bool
	CollisionEnter(other Collider)
	CollisionExit(other Collider)
}

const FatAABBFactor float64 = 0.4

type ShapeCollider interface {
	IntersectSphere(other *SphereCollider)
}

type SphereCollider struct {
	ShapeCollider
	sphere   geo.Sphere
	position vector3.Vector3
}

func (s *SphereCollider) AABB(scale float64, fatter float64) *AABB {
	if fatter == 0 {
		return NewAABB(s.position, vector3.MulScalar(vector3.One, 2*s.sphere.Radius))
	} else {
		size := vector3.MulScalar(vector3.MulScalar(vector3.One, 2*(s.sphere.Radius*scale)), 1+fatter)
		return NewAABB(s.position, size)
	}
}

func (s *SphereCollider) SetPosition(position vector3.Vector3) {
	s.position = position
}

func (s *SphereCollider) IntersectSphere(other *SphereCollider) bool {
	magnitude := s.position.Sub(other.position).SqrMagnitude()
	sr := s.sphere.Radius
	or := other.sphere.Radius
	square := math64.Square(sr + or)
	return magnitude <= square
}
