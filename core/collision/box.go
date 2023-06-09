package collision

import (
	"game-engine/core/geo"
)

type BoxCollider struct {
	ShapeCollider
	box geo.Box
}

func (s *BoxCollider) Type() geo.ShapeType {
	return s.Type()
}

// expend size fatter
func (s *BoxCollider) AABB(fatter float64) *AABB {
	aabb := NewAABB(s.Center(), s.box.Size)
	if fatter > 0 {
		aabb.Expend(1.0 + fatter)
	}
	return aabb
}

func (s *BoxCollider) IntersectSphere(other *SphereCollider) bool {
	return true
}

func (s *BoxCollider) IntersectBox(other *BoxCollider) bool {
	return true
}

//func (a *BoxCollider) ContainsPoint(point vector3.Vector3) bool {
//	return (point.X >= a.min.X && point.X <= a.max.X) &&
//		(point.Y >= a.min.Y && point.Y <= a.max.Y) &&
//		(point.Z >= a.min.Z && point.Z <= a.max.Z)
//}
//
//func (a *BoxCollider) Contains(other *AABB) bool {
//	return a.ContainsPoint(other.min) && a.ContainsPoint(other.max)
//}
//
//func (a *BoxCollider) Intersect(other *AABB) bool {
//	return (a.min.X <= other.max.X && a.max.X >= other.min.X) &&
//		(a.min.Y <= other.max.Y && a.max.Y >= other.min.Y) &&
//		(a.min.Z <= other.max.Z && a.max.Z >= other.min.Z)
//}
