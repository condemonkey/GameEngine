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
	extents := size.MulScalar(0.5)
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

func (a *AABB) Raycast(ray *Ray, maxDistance float64, hit *RaycastHit) bool {
	tmin := -math.MaxFloat64
	tmax := math.MaxFloat64

	const epsilon float64 = 0.00001

	point := ray.Point(maxDistance)
	direction := ray.Direction()

	for i := 0; i < 3; i++ {
		// If the ray is parallel to the slab
		if math.Abs(direction.At(i)) < epsilon {
			// If origin of the ray is not inside the slab, no hit
			if point.At(i) < a.min.At(i) || point.At(i) > a.max.At(i) {
				return false
			}
		} else {
			rayDirectionInverse := 1 / direction.At(i)
			t1 := (a.min.At(i) - point.At(i)) * rayDirectionInverse
			t2 := (a.max.At(i) - point.At(i)) * rayDirectionInverse
			if t1 > t2 {
				// Swap t1 and t2
				tTemp := t2
				t2 = t1
				t1 = tTemp
			}
			tmin = math.Max(tmin, t1)
			tmax = math.Min(tmax, t2)

			// Exit with no collision
			if tmin > tmax {
				return false
			}
		}
	}

	// Compute the hit point
	//hit = point + tMin*rayDirection
	return true

	//if maxDistance > 0 {
	//	tmax = math.Inf(1)
	//}
	//e := a.extents
	//o := vector3.Sub(ray.Origin(), a.center)
	//d := ray.Direction()
	//invD := vector3.DivScalar(1.0, d)
	//
	//var t [6]float64
	//for i := 0; i < 3; i++ {
	//	t[2*i] = -(e.At(i) + o.At(i)) * invD.At(i)
	//	t[2*i+1] = (e.At(i) - o.At(i)) * invD.At(i)
	//}
	//
	//tmin = math64.Max(math.Min(t[0], t[1]), math.Min(t[2], t[3]), math.Min(t[4], t[5]))
	//tmax = math64.Min(math.Max(t[0], t[1]), math.Max(t[2], t[3]), math.Max(t[4], t[5]))
	//if tmax < 0 || tmin > tmax {
	//	return false
	//}
	//if tmin < 0 {
	//	tmin = tmax
	//}
	//hit = NewRaycastHit(nil, tmin, vector3.Sum(o, d.MulScalar(tmin)), vector3.Zero)
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
