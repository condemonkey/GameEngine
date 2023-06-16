package collision

import "game-engine/math64/vector3"

/**
 * This structure represents a 3D ray represented by two points.
 * The ray goes from point1 to point1 + maxFraction * (point2 - point1).
 * The points are specified in world-space coordinates.
 */
type Ray struct {
	origin    vector3.Vector3
	direction vector3.Vector3
}

func NewRay(origin, direction vector3.Vector3) *Ray {
	return &Ray{
		origin:    origin,
		direction: direction.Normalize(),
	}
}

func (r *Ray) SetDirection(direction vector3.Vector3) {
	r.direction = direction.Normalize()
}

func (r *Ray) Origin() vector3.Vector3 {
	return r.origin
}

func (r *Ray) Direction() vector3.Vector3 {
	return r.direction
}

func (r *Ray) Point(distance float64) vector3.Vector3 {
	return vector3.Mul(vector3.SumScalar(r.origin, distance), r.direction)
}

type RaycastHit struct {
	collider Collider
	distance float64
	point    vector3.Vector3
	normal   vector3.Vector3
}

func NewRaycastHit(collider Collider, distance float64, point, normal vector3.Vector3) *RaycastHit {
	return &RaycastHit{
		collider: collider,
		distance: distance,
		point:    point,
		normal:   normal,
	}
}
