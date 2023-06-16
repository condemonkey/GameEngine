package collision

import (
	"game-engine/math64/vector3"
	"math/rand"
	"testing"
	"time"
)

func TestOCTreeRange(t *testing.T) {
	size := float64(4000)
	ecount := 4000
	tree := NewOctree(vector3.Vector3{X: size, Y: size, Z: size})
	now := time.Now()
	for i := 0; i < ecount; i++ {
		position := vector3.RandSpherePoint(vector3.Zero, float64(size/2)-50)
		col := NewBoxCollider(position, vector3.Vector3{X: 5, Y: 5, Z: 5})
		if !tree.Insert(col) {
			aabb := col.AABB(0)
			t.Log("fail insert", aabb.Min(), aabb.Max())
		}
	}
	t.Log(time.Now().Sub(now).Milliseconds())
	col := NewBoxCollider(vector3.Zero, vector3.Vector3{X: 5, Y: 5, Z: 5})
	if !tree.Insert(col) {
		aabb := col.AABB(0)
		t.Log("fail insert", aabb.Min(), aabb.Max())
	}

	collisions := 0
	now = time.Now()
	for i := 0; i < ecount; i++ {
		result := tree.GetColliding(NewAABBWithSize(vector3.Zero, vector3.Vector3{X: 20, Y: 20, Z: 20}))
		collisions = len(result)
	}
	t.Log(time.Now().Sub(now).Milliseconds())

	t.Log("collisions", collisions)
	//t.Log(result.(*TestCollider).Center())
}

func TestOCTree(t *testing.T) {
	//min, max := -2000, 2000
	//ymin, ymax := 1, 1
	size := float64(10000)
	ecount := 8000
	rand.Seed(int64(ecount))

	//wg := new(sync.WaitGroup)
	tree := NewOctree(vector3.Vector3{X: size, Y: size, Z: size})
	var colliders []Collider

	now := time.Now()
	for i := 0; i < ecount; i++ {
		position := vector3.RandSpherePoint(vector3.Zero, float64(size/2)-50)
		col := NewBoxCollider(position, vector3.Vector3{X: 5, Y: 5, Z: 5})
		if !tree.Insert(col) {
			aabb := col.AABB(0)
			t.Log("fail insert", aabb.Min(), aabb.Max())
		}
		//t.Log(col.Id())
		colliders = append(colliders, col)
	}
	t.Log(time.Now().Sub(now).Milliseconds())

	now = time.Now()
	for i := 0; i < ecount; i++ {
		position := vector3.RandSpherePoint(vector3.Zero, float64(size/2)-10)
		if !tree.Move(colliders[i], position) {
			aabb := colliders[i].AABB(0)
			t.Log("fail move", aabb.Min(), aabb.Max())
		}
	}
	t.Log(time.Now().Sub(now).Milliseconds())
}
