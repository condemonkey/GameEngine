package collision

import (
	"game-engine/math64"
	"game-engine/math64/vector3"
	"sync/atomic"
)

const FatAABBFactor float64 = 0.4

type Hittable interface {
	OnCollisionEnter(handle Hittable)
}

type ShapeType int

const (
	ShapeSphere  ShapeType = 1
	ShapeCapsule           = 2
	ShapeBox               = 3
	ShapePlane             = 4
)

type Shape interface {
	AABB(expand float64) *AABB
	IntersectSphere(other *SphereCollider) bool
	IntersectBox(other *BoxCollider) bool
	IntersectDistance(maxDistance float64, origin vector3.Vector3) bool
	IntersectRay(ray Ray) bool
	Type() ShapeType
}

func Box2Box(boxA, boxB *BoxCollider) bool {
	// 현재 center 기준 min,max값이 필요 매번 연산 해야함.
	aMin, aMax := boxA.MinMax()
	bMin, bMax := boxB.MinMax()
	return (aMin.X <= bMax.X && aMax.X >= bMin.X) &&
		(aMin.Y <= bMax.Y && aMax.Y >= bMin.Y) &&
		(aMin.Z <= bMax.Z && aMax.Z >= bMin.Z)
}

func Sphere2Sphere(sphereA, sphereB *SphereCollider) bool {
	magnitude := sphereA.Center().Sub(sphereB.collider.Center()).SqrMagnitude()
	sr := sphereA.Radius()
	or := sphereB.Radius()
	square := math64.Square(sr + or)
	return magnitude <= square
}

func Box2Sphere(box *BoxCollider, sphere *SphereCollider) bool {
	sphereCenter := sphere.Center()
	boxCenter := box.Center()
	cX := math64.Clamp(sphereCenter.X, boxCenter.X-box.with/2, boxCenter.X+box.with/2)
	cY := math64.Clamp(sphereCenter.Y, boxCenter.Y-box.height/2, boxCenter.Y+box.height/2)
	cZ := math64.Clamp(sphereCenter.Z, boxCenter.Z-box.depth/2, boxCenter.Z+box.depth/2)
	distanceSquare := (cX-sphereCenter.X)*(cX-sphereCenter.X) + (cY-sphereCenter.Y)*(cY-sphereCenter.Y) + (cZ-sphereCenter.Z)*(cZ-sphereCenter.Z)
	return distanceSquare < sphere.Radius()*sphere.Radius()
}

type Collider interface {
	Id() int
	Shape() Shape
	Intersect(col Collider) bool
	Center() vector3.Vector3
	SetCenter(center vector3.Vector3)
	AABB(expend float64) *AABB
	CollisionEnter(col Collider)
}

type TestCollider struct {
	Collider
	center vector3.Vector3
	shape  Shape
	id     int
}

var (
	idInc int32
)

func newIdc() int {
	return int(atomic.AddInt32(&idInc, 1))
}

func (c *TestCollider) Id() int {
	return c.id
}

func (c *TestCollider) Shape() Shape {
	return c.shape
}

func (c *TestCollider) Intersect(col Collider) bool {
	targetShape := col.Shape()
	switch targetShape.Type() {
	case ShapeSphere:
		return c.Shape().IntersectSphere(targetShape.(*SphereCollider))
	case ShapeBox:
		return c.Shape().IntersectBox(targetShape.(*BoxCollider))
	default:
		panic("invalid collision shape type")
	}
}

func (c *TestCollider) Center() vector3.Vector3 {
	return c.center
}

func (c *TestCollider) SetCenter(center vector3.Vector3) {
	c.center = center
}

func (c *TestCollider) AABB(expend float64) *AABB {
	return c.shape.AABB(expend)
}

func NewSphereCollider(center vector3.Vector3, radius float64) Collider {
	collider := &TestCollider{}
	collider.id = newIdc()
	collider.center = center
	collider.shape = NewSphere(collider, radius)
	return collider
}

func NewBoxCollider(center vector3.Vector3, size vector3.Vector3) Collider {
	collider := &TestCollider{}
	collider.id = newIdc()
	collider.center = center
	collider.shape = NewBox(collider, size)
	return collider
}
