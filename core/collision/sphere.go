package collision

import (
	"game-engine/math64/vector3"
)

type SphereCollider struct {
	Shape
	radius   float64
	collider Collider
}

func NewSphere(collider Collider, radius float64) Shape {
	shape := &SphereCollider{
		radius:   radius,
		collider: collider,
	}
	return shape
}

func (s *SphereCollider) Center() vector3.Vector3 {
	return s.collider.Center()
}

func (s *SphereCollider) Type() ShapeType {
	return ShapeSphere
}

func (s *SphereCollider) Radius() float64 {
	return s.radius //* s.collider.scale.Max() // 일단 스케일 적용은 X
}

// expend size fatter
func (s *SphereCollider) AABB(fatter float64) *AABB {
	if fatter == 0 {
		return NewAABBWithSize(s.Center(), vector3.MulScalar(vector3.One, 2*s.Radius()))
	} else {
		return NewAABBWithSize(s.Center(), vector3.MulScalar(vector3.MulScalar(vector3.One, 2*s.Radius()), 1+fatter))
	}
}

func (s *SphereCollider) IntersectSphere(other *SphereCollider) bool {
	return Sphere2Sphere(s, other)
}

func (s *SphereCollider) IntersectBox(other *BoxCollider) bool {
	return Box2Sphere(other, s)
}

func (s *SphereCollider) IntersectDistance(maxDistance float64, origin vector3.Vector3) bool {
	return vector3.Distance(origin, s.Center())-s.Radius() < maxDistance
}

func (s *SphereCollider) IntersectRay(other Ray) bool {
	return true
}
