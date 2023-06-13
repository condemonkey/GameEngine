package collision

import (
	"game-engine/core/transform"
	"game-engine/math64/vector3"
)

//type ColliderType int
//
//const (
//	ShapeSphere  ColliderType = 1
//	ShapeCapsule              = 2
//	ShapeBox                  = 3
//)
//
//type Collider interface {
//	Type() ColliderType
//	FatAABB() *AABB
//	AABB() *AABB
//	IntersectRay(ray *Ray, maxDistance float64, hit *RaycastHit) bool
//	IntersectShape(other Collider) bool
//	CollisionEnter(other Collider)
//	CollisionExit(other Collider)
//}

const FatAABBFactor float64 = 0.4

type Hittable interface {
	OnHit(handle Hittable)
	Position() vector3.Vector3
	Scale() vector3.Vector3
}

type OnCollisionFunc func(hit *Collider)

//type Collider interface {
//	Intersect(other Collider) bool
//	IntersectSphere(other *Sphere) bool
//	IntersectDistance(maxDistance float64, origin vector3.Vector3) bool
//	AABB(expand float64) *AABB
//	SetPosition(position vector3.Vector3)
//	Center() vector3.Vector3
//	SetScale(scale vector3.Vector3)
//	SetHandle(user Hittable)
//	SetCollisionHandle(handler OnCollisionFunc)
//	CallCollisionHandle(hit Collider)
//	Handle() Hittable
//	Type() geo.ShapeType
//}

type ShapeType int

const (
	ShapeSphere  ShapeType = 1
	ShapeCapsule           = 2
	ShapeBox               = 3
	ShapePlane             = 4
)

type ShapeCollider interface {
	AABB(expand float64) *AABB
	IntersectSphere(other *Sphere) bool
	IntersectDistance(maxDistance float64, origin vector3.Vector3) bool
	Type() ShapeType
	Center() vector3.Vector3
}

type Collider struct {
	transform *transform.Transform // 부모 transform
	center    vector3.Vector3      // 부모기준 센터
	handle    Hittable
	handler   OnCollisionFunc
	shape     ShapeCollider
}

func NewSphereCollider(radius float64) *Collider {
	collider := &Collider{
		transform: transform.NewTransform(),
	}
	shape := &Sphere{
		radius:   radius,
		collider: collider,
	}
	collider.shape = shape
	return collider
}

//func NewBoxCollider(size vector3.Vector3) *Collider {
//	collider := &Collider{
//		transform: transform.NewTransform(),
//	}
//	shape := &Sphere{
//		collider: collider,
//	}
//	collider.shape = shape
//	return collider
//}

func (s *Collider) Shape() ShapeCollider {
	return s.shape
}

func (s *Collider) IntersectShape(other *Collider) bool {
	hit := false
	switch other.Type() {
	case ShapeSphere:
		hit = s.shape.IntersectSphere(other.shape.(*Sphere))
		break
	default:
		panic("invalid collision shape type")
	}
	return hit
}

func (s *Collider) IntersectDistance(origin vector3.Vector3, distance float64) bool {
	return s.shape.IntersectDistance(distance, origin)
}

func (s *Collider) Center() vector3.Vector3 {
	return s.transform.Position.Add(s.center)
}

func (s *Collider) SetTransform(transform *transform.Transform) {
	s.transform = transform
}

func (s *Collider) SetPosition(position vector3.Vector3) {
	s.transform.Position = position
}

func (s *Collider) SetScale(scale vector3.Vector3) {
	s.transform.Scale = scale
}

func (s *Collider) Handle() Hittable {
	return s.handle
}

func (s *Collider) Type() ShapeType {
	return s.shape.Type()
}

//func (s *Collider) SetHandle(user Hittable) {
//	s.handle = user
//	s.SetPosition(s.handle.Position())
//	s.SetScale(s.handle.Scale())
//}

//func (s *Collider) SetCollisionHandle(handler OnCollisionFunc) {
//	s.handler = handler
//}
//
//func (s *Collider) CallCollisionHandle(hit *Collider) {
//	if s.handler == nil {
//		return
//	}
//	s.handler(hit)
//}

func (s *Collider) AABB(expend float64) *AABB {
	return s.shape.AABB(expend)
}
