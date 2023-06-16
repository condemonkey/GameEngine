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

func NewAABBWithSize(center, size vector3.Vector3) *AABB {
	extents := size.MulScalar(0.5)
	return &AABB{
		center:  center,
		size:    size,
		min:     center.Sub(extents),
		max:     center.Add(extents),
		extents: extents,
	}
}

func NewAABBWithMinMax(min, max vector3.Vector3) *AABB {
	center := max.Add(min).MulScalar(0.5)
	size := max.Sub(center).MulScalar(2)
	return NewAABBWithSize(center, size)
}

//func (a *AABB) Equal(b *AABB) bool {
//	return a.min.Equal(*other.Min) && a.max.Equal(*other.Max)
//}

func (a *AABB) Center() vector3.Vector3 {
	return a.center
}

func (a *AABB) Size() vector3.Vector3 {
	return a.size
}

func (a *AABB) Extents() vector3.Vector3 {
	return a.extents
}

func (a *AABB) Expend(amount float64) *AABB {
	a.size = a.size.Add(vector3.One.MulScalar(amount))
	a.extents = a.size.MulScalar(0.5)
	a.min = a.center.Sub(a.extents)
	a.max = a.center.Add(a.extents)
	return a
}

func (a *AABB) OctSplit() [8]*AABB {
	return [8]*AABB{
		NewAABBWithMinMax(vector3.Vector3{X: a.min.X, Y: a.min.Y, Z: a.min.Z}, vector3.Vector3{X: a.center.X, Y: a.center.Y, Z: a.center.Z}),
		NewAABBWithMinMax(vector3.Vector3{X: a.min.X, Y: a.min.Y, Z: a.center.Z}, vector3.Vector3{X: a.center.X, Y: a.center.Y, Z: a.max.Z}),
		NewAABBWithMinMax(vector3.Vector3{X: a.min.X, Y: a.center.Y, Z: a.min.Z}, vector3.Vector3{X: a.center.X, Y: a.max.Y, Z: a.center.Z}),
		NewAABBWithMinMax(vector3.Vector3{X: a.min.X, Y: a.center.Y, Z: a.center.Z}, vector3.Vector3{X: a.center.X, Y: a.max.Y, Z: a.max.Z}),

		NewAABBWithMinMax(vector3.Vector3{X: a.center.X, Y: a.min.Y, Z: a.min.Z}, vector3.Vector3{X: a.max.X, Y: a.center.Y, Z: a.center.Z}),
		NewAABBWithMinMax(vector3.Vector3{X: a.center.X, Y: a.min.Y, Z: a.center.Z}, vector3.Vector3{X: a.max.X, Y: a.center.Y, Z: a.max.Z}),
		NewAABBWithMinMax(vector3.Vector3{X: a.center.X, Y: a.center.Y, Z: a.min.Z}, vector3.Vector3{X: a.max.X, Y: a.max.Y, Z: a.center.Z}),
		NewAABBWithMinMax(vector3.Vector3{X: a.center.X, Y: a.center.Y, Z: a.center.Z}, vector3.Vector3{X: a.max.X, Y: a.max.Y, Z: a.max.Z}),
	}
}

func (a *AABB) Min() vector3.Vector3 {
	return a.min
}

func (a *AABB) Max() vector3.Vector3 {
	return a.max
}

func (a *AABB) With() float64 {
	return a.max.X - a.min.X
}

func (a *AABB) Height() float64 {
	return a.max.Y - a.min.Y
}

func (a *AABB) Depth() float64 {
	return a.max.Z - a.min.Z
}

func (a *AABB) ContainsPoint(point vector3.Vector3) bool {
	return (point.X >= a.min.X && point.X <= a.max.X) &&
		(point.Y >= a.min.Y && point.Y <= a.max.Y) &&
		(point.Z >= a.min.Z && point.Z <= a.max.Z)
}

func (a *AABB) Contains(b *AABB) bool {
	return a.ContainsPoint(b.min) && a.ContainsPoint(b.max)
}

//func (a *AABB) Contains2(v vector3.Vector3) bool {
//	return (a.min.X <= v.X && v.X <= a.max.X) &&
//		(a.min.Y <= v.Y && v.Y <= a.max.Y) &&
//		(a.min.Z <= v.Z && v.Z <= a.max.Z)
//}
//
//func (a *AABB) Fit(b *AABB) bool {
//	return b.Contains2(a.max) && a.Contains2(a.min)
//}

func (a *AABB) Intersect(b *AABB) bool {
	return !(a.max.X < b.min.X || b.max.X < a.min.X || a.max.Y < b.min.Y || b.max.Y < a.min.Y || a.max.Z < b.min.Z || b.max.Z < a.min.Z)
	//return (a.min.X <= b.max.X && a.max.X >= b.min.X) &&
	//	(a.min.Y <= b.max.Y && a.max.Y >= b.min.Y) &&
	//	(a.min.Z <= b.max.Z && a.max.Z >= b.min.Z)
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

func (a *AABB) IntersectRay(ray Ray) bool {
	dirfrac := vector3.Vector3{}
	dirfrac.X = 1.0 / ray.direction.X
	dirfrac.Y = 1.0 / ray.direction.Y
	dirfrac.Z = 1.0 / ray.direction.Z
	// lb is the corner of AABB with minimal coordinates - left bottom, rt is maximal corner
	// r.org is origin of ray
	t1 := (a.min.X - ray.origin.X) * dirfrac.X
	t2 := (a.max.X - ray.origin.X) * dirfrac.X
	t3 := (a.min.Y - ray.origin.Y) * dirfrac.Y
	t4 := (a.max.Y - ray.origin.Y) * dirfrac.Y
	t5 := (a.min.Z - ray.origin.Z) * dirfrac.Z
	t6 := (a.max.Z - ray.origin.Z) * dirfrac.Z

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

func AABBExpand(a *AABB, amount float64) {
	a.size = vector3.Sum(a.size, vector3.MulScalar(vector3.One, amount))
	a.extents = vector3.MulScalar(a.size, 0.5)
	a.min = vector3.Sub(a.center, a.extents)
	a.max = vector3.Sum(a.center, a.extents)
}

func AABBIntersect(a *AABB, b *AABB) bool {
	return (a.min.X <= b.max.X && a.max.X >= b.min.X) &&
		(a.min.Y <= b.max.Y && a.max.Y >= b.min.Y) &&
		(a.min.Z <= b.max.Z && a.max.Z >= b.min.Z)
}

//func AABBOverlap(a *AABB, b *AABB) bool {
//	if a.max.X < b.min.X || b.max.X < a.min.X {
//		return false
//	}
//	if a.max.Y < b.min.Y || b.max.Y < a.min.Y {
//		return false
//	}
//	if a.max.Z < b.min.Z || b.max.Z < a.min.Z {
//		return false
//	}
//	return true
//}

func AABBUnion(a, b *AABB) *AABB {
	max := vector3.Max(a.max, b.max)
	min := vector3.Min(a.min, b.min)
	center := max.Add(min).MulScalar(0.5)
	size := max.Sub(center).MulScalar(2)
	return &AABB{
		min:    min,
		max:    max,
		size:   size,
		center: center,
	}
}
