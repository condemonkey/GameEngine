package core

import (
	"game-engine/core/collision"
	"game-engine/math64/vector3"
)

const FatAABBFactor float64 = 0.4

type Collider struct {
	BaseComponent
	UpdatableComponent
	collision.Collider
	center   vector3.Vector3
	collider Shape
}

func NewSphereCollider() *Collider {
	return &Collider{
		BaseComponent: NewBaseComponent(),
		collider:      NewSphere(),
	}
}

func (c *Collider) Update(dt int) {
}

func (c *Collider) Awake() {
	c.collider.SetTransform(c.Transform())
}

func (c *Collider) Raycast(ray collision.Ray, maxDistance float64) *collision.RaycastHit {
	//https://noti.st/eiaserinnys/jCpSbp/slides
	return nil
	//return c.collider.Raycast(ray, maxDistance)
}

func (c *Collider) Center() vector3.Vector3 {
	return vector3.Sum(c.Transform().Position(), c.center)
}

func (c *Collider) Type() collision.ColliderType {
	return c.collider.Type()
}

func (c *Collider) InternalShape() Shape {
	return c.collider
}

func (c *Collider) FatAABB() *collision.AABB {
	return c.collider.FatAABB()
}

func (c *Collider) AABB() *collision.AABB {
	return c.collider.AABB()
}

type Shape interface {
	collision.Collider
	Center() vector3.Vector3
	SetTransform(transform *Transform)
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
func NewSphere() Shape {
	return &Sphere{
		ShapeCollider: ShapeCollider{
			center:    vector3.Zero,
			transform: nil,
		},
		radius: 0.5,
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

func (s *Sphere) Raycast(ray *collision.Ray, maxDistance float64, hit *collision.RaycastHit) bool {
	return nil
}

//type BoxCollider struct {
//	Collider
//}
//
//type CapsuleCollider struct {
//	Collider
//}
