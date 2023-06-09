package collision

import (
	"game-engine/core/geo"
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

type OnCollisionFunc func(hit Collider)

type Collider interface {
	Intersect(other Collider) bool
	IntersectSphere(other *SphereCollider) bool
	AABB(expand float64) *AABB
	SetPosition(position vector3.Vector3)
	Center() vector3.Vector3
	SetScale(scale vector3.Vector3)
	SetHandle(user Hittable)
	SetCollisionHandle(handler OnCollisionFunc)
	CallCollisionHandle(hit Collider)
	Handle() Hittable
	Type() geo.ShapeType
}

func NewCollider(shape geo.Shape) Collider {
	switch shape.Type() {
	case geo.ShapeSphere:
		return &SphereCollider{
			ShapeCollider: &ShapeCollider{
				scale:       vector3.One,
				parentScale: vector3.One,
			},
			sphere: shape.(geo.Sphere),
		}
	}
	panic("")
}

type ShapeCollider struct {
	Collider
	parentPosition vector3.Vector3 // 부모 포지션
	parentScale    vector3.Vector3 // 부모 스케일
	scale          vector3.Vector3
	center         vector3.Vector3
	handle         Hittable
	handler        OnCollisionFunc
}

func (s *ShapeCollider) Center() vector3.Vector3 {
	return s.parentPosition.Add(s.center)
}

func (s *ShapeCollider) SetPosition(position vector3.Vector3) {
	s.parentPosition = position
}

func (s *ShapeCollider) SetScale(scale vector3.Vector3) {
	s.parentScale = scale
}

func (s *ShapeCollider) Handle() Hittable {
	return s.handle
}

func (s *ShapeCollider) SetHandle(user Hittable) {
	s.handle = user
	s.SetPosition(s.handle.Position())
	s.SetScale(s.handle.Scale())
}

func (s *ShapeCollider) SetCollisionHandle(handler OnCollisionFunc) {
	s.handler = handler
}

func (s *ShapeCollider) CallCollisionHandle(hit Collider) {
	if s.handler == nil {
		return
	}
	s.handler(hit)
}

func (s *ShapeCollider) AABB(expend float64) *AABB {
	return s.AABB(expend)
}
