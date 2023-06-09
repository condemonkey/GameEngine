package core

import (
	"game-engine/core/collision"
	"game-engine/math64"
	"game-engine/math64/vector3"
)

type Shape interface {
	Center() vector3.Vector3
	SetTransform(transform *Transform)
	Type() collision.ColliderType
	FatAABB() *collision.AABB
	AABB() *collision.AABB

	// collide
	IntersectRay(ray *collision.Ray, maxDistance float64, hit *collision.RaycastHit) bool
	IntersectSphere(other *Sphere) bool
}

type ShapeCollider struct {
	Shape
	transform *Transform
	center    vector3.Vector3
}

func (s *ShapeCollider) Center() vector3.Vector3 {
	return s.transform.Position().Add(s.center)
}

func (s *ShapeCollider) SetTransform(transform *Transform) {
	s.transform = transform
}

// internal shape
func NewSphere(radius float64) Shape {
	return &Sphere{
		ShapeCollider: ShapeCollider{
			center: vector3.Zero,
		},
		radius: radius,
	}
}

type Sphere struct {
	ShapeCollider
	radius float64
}

func (s *Sphere) Radius() float64 {
	return s.radius * s.transform.Scale().Max()
}

func (s *Sphere) Type() collision.ColliderType {
	return collision.ShapeSphere
}

func (s *Sphere) FatAABB() *collision.AABB {
	size := vector3.MulScalar(vector3.MulScalar(vector3.One, 2*s.Radius()), 1+FatAABBFactor)
	return collision.NewAABB(s.Center(), size)
}

func (s *Sphere) AABB() *collision.AABB {
	return collision.NewAABB(s.Center(), vector3.MulScalar(vector3.One, 2*s.Radius()))
}

func (s *Sphere) IntersectRay(ray *collision.Ray, maxDistance float64, hit *collision.RaycastHit) bool {
	return true
}

func (s *Sphere) IntersectSphere(other *Sphere) bool {
	// sphere vs sphere
	magnitude := s.Center().Sub(other.Center()).SqrMagnitude()
	sr := s.Radius()
	or := other.Radius()
	square := math64.Square(sr + or)
	return magnitude <= square
}

//type BoxCollider struct {
//	Collider
//}
//
//type CapsuleCollider struct {
//	Collider
//}
