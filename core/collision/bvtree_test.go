package collision

import (
	"game-engine/math64"
	"game-engine/math64/vector3"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestCollision(t *testing.T) {
	box := NewBoxCollider(vector3.Vector3{X: 0, Y: 0, Z: 0}, vector3.Vector3{X: 10, Y: 10, Z: 10})
	sphere := NewSphereCollider(vector3.Vector3{X: 10, Y: 10, Z: 10}, 10)
	//collider2 := NewSphereCollider(vector3.Vector3{X: 0, Y: 0, Z: 0}, 10)

	start := time.Now()
	for j := 0; j < 100000; j++ {
		box.Intersect(sphere)
	}
	t.Log(time.Since(start).Milliseconds())
}

func TestTreeRelocate(t *testing.T) {
	ecount := 8000
	rand.Seed(int64(ecount))

	wg := new(sync.WaitGroup)
	tree := NewBVTree(wg)

	var spheres []Collider
	for i := 0; i < ecount; i++ {
		position := vector3.RandSpherePoint(vector3.Zero, float64(ecount)-1)
		size := vector3.Vector3{X: 5, Y: 5, Z: 5}
		spheres = append(spheres, NewBoxCollider(position, size))
	}

	now := time.Now()
	for _, entity := range spheres {
		tree.AddCollider(entity)
	}
	t.Log(time.Since(now).Milliseconds())

	now = time.Now()
	for _, entity := range spheres {
		position := vector3.RandSpherePoint(vector3.Zero, float64(ecount)-1)
		entity.(*TestCollider).SetCenter(position)
		tree.RelocateCollider(entity)
	}
	t.Log(time.Since(now).Milliseconds())

}

func NewTestRandSphereCollider(min, max, minY, maxY int, radius float64) Collider {
	x := float64(math64.RandInt(min, max))
	y := float64(math64.RandInt(minY, maxY))
	z := float64(math64.RandInt(min, max))
	return NewSphereCollider(vector3.Vector3{X: x, Y: y, Z: z}, radius)
}

func NewTestRandBoxCollider(min, max, minY, maxY int, size float64) Collider {
	x := float64(math64.RandInt(min, max))
	y := float64(math64.RandInt(minY, maxY))
	z := float64(math64.RandInt(min, max))
	return NewBoxCollider(vector3.Vector3{X: x, Y: y, Z: z}, vector3.Vector3{X: size, Y: size, Z: size})
}

func TestTree(t *testing.T) {
	radius := 2000 // 이러면 4000 원형 범위가 생김
	ecount := 8000
	rand.Seed(int64(ecount))

	wg := new(sync.WaitGroup)
	tree := NewBVTree(wg)

	var spheres []Collider
	for i := 0; i < ecount; i++ {
		position := vector3.RandSpherePoint(vector3.Zero, float64(radius))
		size := vector3.Vector3{X: 20, Y: 20, Z: 20}
		spheres = append(spheres, NewBoxCollider(position, size))
	}

	distance := float64(600)
	position := vector3.RandSpherePoint(vector3.Zero, float64(radius))

	collider := NewSphereCollider(position, distance)
	tree.AddCollider(collider)

	now := time.Now()
	for _, entity := range spheres {
		tree.AddCollider(entity)
	}
	t.Log(time.Now().Sub(now).Milliseconds())

	//origin := vector3.Vector3{X: x, Y: y, Z: z}

	loop := 20
	//for i := 0; i < loop; i++ {
	//	count := 0
	//	now = time.Now()
	//	for j := 0; j < ecount; j++ {
	//		tree.CollidingBox(origin, distance, func(colliders []Collider) {
	//			count = len(colliders)
	//		})
	//	}
	//	tree.WaitGroup()
	//	t.Log("detect time1", time.Now().Sub(now).Milliseconds(), "collide", count)
	//	time.Sleep(time.Millisecond * 30)
	//}
	for i := 0; i < loop; i++ {
		count := 0
		now = time.Now()
		for j := 0; j < ecount; j++ {
			tree.CollidingSphereAsync(position, distance, func(colliders []Collider) {
				count = len(colliders)
			})
		}
		tree.WaitGroup()
		t.Log("detect time2", time.Now().Sub(now).Milliseconds(), "collide", count)
		time.Sleep(time.Millisecond * 30)
	}

	tree.Snapshot()
}
