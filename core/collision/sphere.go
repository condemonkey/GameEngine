package collision

import (
	"game-engine/math64"
	"game-engine/math64/vector3"
)

type Sphere struct {
	ShapeCollider
	radius   float64
	collider *Collider
}

func (s *Sphere) Center() vector3.Vector3 {
	return s.collider.Center()
}

func (s *Sphere) Type() ShapeType {
	return ShapeSphere
}

func (s *Sphere) Radius() float64 {
	return s.radius //* s.collider.scale.Max() // 일단 스케일 적용은 X
}

// expend size fatter
func (s *Sphere) AABB(fatter float64) *AABB {
	if fatter == 0 {
		return NewAABB(s.Center(), vector3.MulScalar(vector3.One, 2*s.Radius()))
	} else {
		return NewAABB(s.Center(), vector3.MulScalar(vector3.MulScalar(vector3.One, 2*s.Radius()), 1+fatter))
	}
}

func (s *Sphere) IntersectSphere(other *Sphere) bool {
	magnitude := s.Center().Sub(other.collider.Center()).SqrMagnitude()
	sr := s.Radius()
	or := other.Radius()
	square := math64.Square(sr + or)
	return magnitude <= square
}

func (s *Sphere) IntersectDistance(maxDistance float64, origin vector3.Vector3) bool {
	return vector3.Distance(origin, s.Center())-s.Radius() < maxDistance
}
