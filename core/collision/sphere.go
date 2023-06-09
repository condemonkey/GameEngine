package collision

import (
	"game-engine/core/geo"
	"game-engine/math64"
	"game-engine/math64/vector3"
)

type SphereCollider struct {
	*ShapeCollider
	sphere geo.Sphere
}

func (s *SphereCollider) Intersect(a Collider) bool {
	hit := false
	switch a.Type() {
	case geo.ShapeSphere:
		hit = s.IntersectSphere(a.(*SphereCollider))
		break
	default:
		panic("invalid collision shape type")
	}

	//if hit && s.handle != nil {
	//	s.handle.OnHit(a.Handle())
	//}

	if hit {
		s.CallCollisionHandle(a)
		a.CallCollisionHandle(s)
	}

	return hit
}

func (s *SphereCollider) Type() geo.ShapeType {
	return s.sphere.Type()
}

func (s *SphereCollider) Radius() float64 {
	return s.sphere.Radius //* s.scale.Max() // 일단 스케일 적용은 X
}

// expend size fatter
func (s *SphereCollider) AABB(fatter float64) *AABB {
	if fatter == 0 {
		return NewAABB(s.Center(), vector3.MulScalar(vector3.One, 2*s.Radius()))
	} else {
		return NewAABB(s.Center(), vector3.MulScalar(vector3.MulScalar(vector3.One, 2*s.Radius()), 1+fatter))
	}
}

func (s *SphereCollider) IntersectSphere(other *SphereCollider) bool {
	magnitude := s.Center().Sub(other.Center()).SqrMagnitude()
	sr := s.Radius()
	or := other.Radius()
	square := math64.Square(sr + or)
	return magnitude <= square
}
