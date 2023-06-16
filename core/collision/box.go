package collision

import (
	"game-engine/math64/vector3"
	"math"
)

type BoxCollider struct {
	Shape
	collider Collider
	size     vector3.Vector3
	with     float64
	height   float64
	depth    float64
}

func NewBox(collider Collider, size vector3.Vector3) Shape {
	shape := &BoxCollider{
		collider: collider,
		size:     size,
	}
	shape.with, shape.height, shape.depth = shape.WithHeightDepth()
	return shape
}

func (s *BoxCollider) Type() ShapeType {
	return ShapeBox
}

func (s *BoxCollider) Center() vector3.Vector3 {
	return s.collider.Center()
}

// expend size fatter
func (s *BoxCollider) AABB(fatter float64) *AABB {
	aabb := NewAABBWithSize(s.collider.Center(), s.size)
	if fatter > 0 {
		aabb.Expend(1.0 + fatter)
	}
	return aabb
}

func (s *BoxCollider) MinMax() (min, max vector3.Vector3) {
	extents := s.size.MulScalar(0.5)
	center := s.collider.Center()
	min = center.Sub(extents)
	max = center.Add(extents)
	return min, max
}

func (s *BoxCollider) WithHeightDepth() (with, height, depth float64) {
	min, max := s.MinMax()
	with = max.X - min.X
	height = max.Y - min.Y
	depth = max.Z - min.Z
	return with, height, depth
}

func (s *BoxCollider) IntersectSphere(other *SphereCollider) bool {
	return Box2Sphere(s, other)
}

func (s *BoxCollider) IntersectBox(other *BoxCollider) bool {
	return Box2Box(s, other)
}

func (s *BoxCollider) IntersectRay(ray Ray) bool {
	min, max := s.MinMax()
	dirfrac := vector3.Vector3{}
	dirfrac.X = 1.0 / ray.direction.X
	dirfrac.Y = 1.0 / ray.direction.Y
	dirfrac.Z = 1.0 / ray.direction.Z
	// lb is the corner of AABB with minimal coordinates - left bottom, rt is maximal corner
	// r.org is origin of ray
	t1 := (min.X - ray.origin.X) * dirfrac.X
	t2 := (max.X - ray.origin.X) * dirfrac.X
	t3 := (min.Y - ray.origin.Y) * dirfrac.Y
	t4 := (max.Y - ray.origin.Y) * dirfrac.Y
	t5 := (min.Z - ray.origin.Z) * dirfrac.Z
	t6 := (max.Z - ray.origin.Z) * dirfrac.Z

	tmin := math.Max(math.Max(math.Min(t1, t2), math.Min(t3, t4)), math.Min(t5, t6))
	tmax := math.Min(math.Min(math.Max(t1, t2), math.Max(t3, t4)), math.Max(t5, t6))

	// if tmax < 0, ray (line) is intersecting AABB, but whole AABB is behing us
	if tmax < 0 {
		return false
	}

	// if tmin > tmax, ray doesn't intersect AABB
	if tmin > tmax {
		return false
	}

	return true
}
