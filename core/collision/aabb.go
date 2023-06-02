package collision

import (
	"game-engine/math64/vector3"
	"math"
)

type AABB struct {
	center  vector3.Vector3
	size    vector3.Vector3
	extents vector3.Vector3
	min     vector3.Vector3
	max     vector3.Vector3
}

func NewAABB(center, size vector3.Vector3) *AABB {
	extents := size.MulScala(0.5)
	return &AABB{
		center:  center,
		size:    size,
		min:     center.Sub(extents),
		max:     center.Add(extents),
		extents: extents,
	}
}

func (a *AABB) Center() vector3.Vector3 {
	return a.center
}

func (a *AABB) Size() vector3.Vector3 {
	return a.size
}

func (a *AABB) Extents() vector3.Vector3 {
	return a.extents
}

func (a *AABB) Min() vector3.Vector3 {
	return a.min
}

func (a *AABB) Max() vector3.Vector3 {
	return a.max
}

func (a *AABB) ContainsPoint(point vector3.Vector3) bool {
	return (point.X >= a.min.X && point.X <= a.max.X) &&
		(point.Y >= a.min.Y && point.Y <= a.max.Y) &&
		(point.Z >= a.min.Z && point.Z <= a.max.Z)
}

func (a *AABB) Contains(other *AABB) bool {
	return a.ContainsPoint(other.min) && a.ContainsPoint(other.max)
}

func (a *AABB) Intersect(other *AABB) bool {
	return (a.min.X <= other.max.X && a.max.X >= other.min.X) &&
		(a.min.Y <= other.max.Y && a.max.Y >= other.min.Y) &&
		(a.min.Z <= other.max.Z && a.max.Z >= other.min.Z)
}

func (a *AABB) SurfaceArea() float64 {
	return 2 * (a.size.X*a.size.Y + a.size.Y*a.size.Z + a.size.X*a.size.Z)
}

func (a *AABB) EncapsulatePoint(point vector3.Vector3) {
	a.min.X = math.Min(a.min.X, point.X)
	a.min.Y = math.Min(a.min.Y, point.Y)
	a.min.Z = math.Min(a.min.Z, point.Z)
	a.max.X = math.Max(a.max.X, point.X)
	a.max.Y = math.Max(a.max.Y, point.Y)
	a.max.Z = math.Max(a.max.Z, point.Z)
}

func AABBExpand(a *AABB, amount float64) {
	a.size = vector3.Sum(a.size, vector3.Scalef(vector3.One, amount))
	a.extents = vector3.Scalef(a.size, 0.5)
	a.min = vector3.Sub(a.center, a.extents)
	a.max = vector3.Sum(a.center, a.extents)
}

func AABBIntersect(a *AABB, b *AABB) bool {
	return (a.min.X <= b.max.X && a.max.X >= b.min.X) &&
		(a.min.Y <= b.max.Y && a.max.Y >= b.min.Y) &&
		(a.min.Z <= b.max.Z && a.max.Z >= b.min.Z)
}

func AABBEncapsulate(a, b *AABB) *AABB {
	max := vector3.Max(a.max, b.max)
	min := vector3.Min(a.min, b.min)
	center := max.Add(min).MulScala(0.5)
	size := max.Sub(center).MulScala(2)
	return &AABB{
		min:    min,
		max:    max,
		size:   size,
		center: center,
	}
}
