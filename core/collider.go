package core

import (
	"game-engine/core/collision"
	"game-engine/math64/vector3"
)

type CollisionHandler interface {
	CollisionEnter(collider *Collider)
	CollisionExit(collider *Collider)
}

const FatAABBFactor float64 = 0.4

type Collider struct {
	BaseComponent
	collision.Collider
	center   vector3.Vector3
	shape    Shape
	handlers []CollisionHandler
}

func NewSphereCollider() *Collider {
	return &Collider{
		BaseComponent: NewBaseComponent(),
		shape:         NewSphere(0.5),
	}
}

func NewShapeCollider(shape Shape) *Collider {
	return &Collider{
		BaseComponent: NewBaseComponent(),
		shape:         shape,
	}
}

func (c *Collider) Awake() {
	c.SetTransform(c.Transform())
}

func (c *Collider) SetTransform(transform *Transform) {
	c.shape.SetTransform(transform)
}

func (c *Collider) IntersectRay(ray *collision.Ray, maxDistance float64, hit *collision.RaycastHit) bool {
	//https://noti.st/eiaserinnys/jCpSbp/slides
	return true
	//return c.collider.Raycast(ray, maxDistance)
}

func (c *Collider) IntersectShape(other collision.Collider) bool {
	shape := other.(*Collider).Shape()
	switch other.Type() {
	case collision.ShapeSphere:
		return c.shape.IntersectSphere(shape.(*Sphere))
	}
	panic("invalid collision shape type")
}

func (c *Collider) Center() vector3.Vector3 {
	return vector3.Sum(c.Transform().Position(), c.center)
}

func (c *Collider) Type() collision.ColliderType {
	return c.shape.Type()
}

func (c *Collider) Shape() Shape {
	return c.shape
}

func (c *Collider) FatAABB() *collision.AABB {
	return c.shape.FatAABB()
}

func (c *Collider) AABB() *collision.AABB {
	return c.shape.AABB()
}

func (c *Collider) AddCollisionHandler(handler CollisionHandler) {
	c.handlers = append(c.handlers, handler)
}

func (c *Collider) RemoveCollisionHandler(handler CollisionHandler) {
}

func (c *Collider) CollisionEnter(other collision.Collider) {
	for _, handler := range c.handlers {
		handler.CollisionEnter(other.(*Collider))
	}
}

func (c *Collider) CollisionExit(other collision.Collider) {
	for _, handler := range c.handlers {
		handler.CollisionExit(other.(*Collider))
	}
}
