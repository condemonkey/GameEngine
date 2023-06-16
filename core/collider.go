package core

import (
	"game-engine/core/collision"
	"game-engine/math64/vector3"
)

type CollisionHandler interface {
	OnCollisionEnter(entity *Entity)
}

type Collider struct {
	BaseComponent
	collision.Collider
	center   vector3.Vector3
	shape    collision.Shape
	handlers []CollisionHandler
}

func NewSphereCollider(radius float64) *Collider {
	collider := &Collider{}
	collider.BaseComponent = NewBaseComponent()
	collider.shape = collision.NewSphere(collider, radius)
	return collider
}

func NewBoxCollider(size vector3.Vector3) *Collider {
	collider := &Collider{}
	collider.BaseComponent = NewBaseComponent()
	collider.shape = collision.NewBox(collider, size)
	return collider
}

func (c *Collider) Id() int {
	return c.Entity().Id()
}

func (c *Collider) Awake() {

}

func (c *Collider) Shape() collision.Shape {
	return c.shape
}

func (c *Collider) Intersect(col collision.Collider) bool {
	targetShape := col.Shape()
	switch targetShape.Type() {
	case collision.ShapeSphere:
		return c.Shape().IntersectSphere(targetShape.(*collision.SphereCollider))
	case collision.ShapeBox:
		return c.Shape().IntersectBox(targetShape.(*collision.BoxCollider))
	default:
		panic("invalid collision shape type")
	}
}

func (c *Collider) Center() vector3.Vector3 {
	return c.Transform().Position().Add(c.center)
}

func (c *Collider) AABB(expend float64) *collision.AABB {
	return c.shape.AABB(expend)
}

func (c *Collider) CollisionEnter(col collision.Collider) {
	entity := c.ConvertEntity(col)
	entity.Collider()
}

func (c *Collider) ConvertEntity(col collision.Collider) *Entity {
	return col.(*Collider).Entity()
}

func ConvertCollider(col collision.Collider) *Entity {
	return col.(*Collider).Entity()
}

//	func (p *Collider) CollisionCollider() *collision.Collider {
//		return p.Collisionable.(*collision.Collider)
//	}

//
//func (p *Collider) SubscribeCollision(collision CollisionHandler) {
//	p.listeners = append(p.listeners, collision)
//}
//
//func (p *Collider) UnsubscribeCollision(collision CollisionHandler) {
//}
