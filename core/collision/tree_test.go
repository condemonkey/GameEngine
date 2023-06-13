package collision

import (
	"fmt"
	"game-engine/math64"
	"game-engine/math64/vector3"
	"math/rand"
	"testing"
	"time"
)

func TestCollision(t *testing.T) {
	collider := NewSphereCollider(0.5)
	collider.SetPosition(vector3.Vector3{X: 0, Y: 0, Z: 0})

	collider2 := NewSphereCollider(0.5)
	collider2.SetPosition(vector3.Vector3{X: 0, Y: 0, Z: 0})

	fmt.Println(collider.IntersectShape(collider2))
	fmt.Println(collider2.IntersectDistance(vector3.Vector3{-10, 10, 10}, 10))
}

func TestTreeUpdatePosition(t *testing.T) {
	//tree := NewBVTree()
	//collider1 := NewCollider(geo.NewSphere(0.5))
	//collider1.SetPosition(vector3.Vector3{X: 0, Y: 0, Z: 0})
	//collider1.SetCollisionHandle(func(hit Collider) {
	//})
	//
	//tree.AddCollider(collider1)
	//
	//collider := NewCollider(geo.NewSphere(0.5))
	//collider.SetPosition(vector3.Vector3{X: 0.6, Y: 0.6, Z: 0.6})
	//collider.SetCollisionHandle(func(hit Collider) {
	//})

	//abs := tree.Query(collider)
	//for _, col := range abs {
	//	col.Intersect(collider)
	//}
	//t.Log(abs)
	//
	//result := collider.Intersect(collider1)
	//t.Log(result)
}

func TestTree(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	min, max := -2000, 2000
	ymin, ymax := 1, 1
	ecount := 5000

	tree := NewBVTree()

	var entities []*Collider
	for i := 0; i < ecount; i++ {
		x := float64(math64.RandInt(min, max))
		y := float64(math64.RandInt(ymin, ymax))
		z := float64(math64.RandInt(min, max))
		collider := NewSphereCollider(1)
		collider.SetPosition(vector3.Vector3{X: x, Y: y, Z: z})
		entities = append(entities, collider)
	}

	distance := float64(300)

	x := float64(math64.RandInt(min, max))
	y := float64(math64.RandInt(ymin, ymax))
	z := float64(math64.RandInt(min, max))
	collider := NewSphereCollider(distance)
	collider.SetPosition(vector3.Vector3{X: x, Y: y, Z: z})
	entities = append(entities, collider)

	now := time.Now()
	for _, entity := range entities {
		tree.AddCollider(entity)
	}
	t.Log(time.Now().Sub(now).Milliseconds())

	origin := vector3.Vector3{X: x, Y: y, Z: z}

	loop := 20
	for i := 0; i < loop; i++ {
		count := 0
		now = time.Now()
		for j := 0; j < ecount; j++ {
			tree.IntersectRangeAsync(origin, distance, func(hits int) {
				count = hits
			})
		}
		tree.WaitGroup()
		t.Log("detect time1", time.Now().Sub(now).Milliseconds(), "collide", count)
		time.Sleep(time.Millisecond * 30)
	}

	for i := 0; i < loop; i++ {
		count := 0
		now = time.Now()
		for j := 0; j < ecount; j++ {
			tree.IntersectRangeCollidersAsync(origin, distance, func(colliders []*Collider) {
				count = len(colliders)
			})
		}
		tree.WaitGroup()
		t.Log("detect time2", time.Now().Sub(now).Milliseconds(), "collide", count)
		time.Sleep(time.Millisecond * 30)
	}

	for i := 0; i < loop; i++ {
		count := 0
		now = time.Now()
		for j := 0; j < ecount; j++ {
			tree.IntersectRangeCollider(origin, distance, func(colliders []*Collider) {
				count = len(colliders)
			})
		}
		t.Log("detect time3", time.Now().Sub(now).Milliseconds(), "collide", count)
		time.Sleep(time.Millisecond * 30)
	}

	//for i := 0; i < loop; i++ {
	//	count := 0
	//	now = time.Now()
	//	for j := 0; j < ecount; j++ {
	//		count = tree.Intersect(collider)
	//	}
	//	t.Log("detect time2", time.Now().Sub(now).Milliseconds(), "collide", count)
	//	time.Sleep(time.Millisecond * 30)
	//}
	//
	//for i := 0; i < loop; i++ {
	//	count := 0
	//	now = time.Now()
	//	for j := 0; j < ecount; j++ {
	//		count = tree.IntersectRange(origin, distance)
	//	}
	//	t.Log("detect time3", time.Now().Sub(now).Milliseconds(), "collide", count)
	//	time.Sleep(time.Millisecond * 30)
	//}

	tree.Snapshot()

	//queue := util.NewQueue[Collider]()
	//for i := 0; i < loop; i++ {
	//	count := 0
	//	now = time.Now()
	//	for j := 0; j < 5000; j++ {
	//		count = tree.IntersectDistance(origin, 20)
	//	}
	//	t.Log("detect time3", time.Now().Sub(now).Milliseconds(), "collide", count)
	//	time.Sleep(time.Millisecond * 30)
	//}
	//
	//count := 0
	//
	//for i := 0; i < loop; i++ {
	//	now = time.Now()
	//	count = 0
	//	for j := 0; j < 5000; j++ {
	//		count = tree.IntersectDistance(origin, 20)
	//	}
	//	t.Log("detect time4", time.Now().Sub(now).Milliseconds(), "collide", count)
	//	time.Sleep(time.Millisecond * 30)
	//}

	//now = time.Now()
	//for _, entity := range entities {
	//	entity.SetPosition(entity.Center().Add(vector3.Vector3{1000, 1000, 1000}))
	//	tree.RelocateCollider(entity)
	//}
	//t.Log(time.Now().Sub(now).Milliseconds())

	//now = time.Now()
	//tree.Update()
	//t.Log(time.Now().Sub(now).Milliseconds())

	//snapshots := Snapshots{}
	//
	//tree.Traverse(func(node *BVTreeNode) {
	//	aabb := node.FatAABB()
	//	radius := float64(0)
	//	if node.IsLeaf() {
	//		radius = node.Collider().Shape().(*Sphere).Radius()
	//	}
	//	snapshots.Snapshots = append(snapshots.Snapshots, Snapshot{
	//		Min:    aabb.Min(),
	//		Max:    aabb.Max(),
	//		Center: aabb.Center(),
	//		Size:   aabb.Size(),
	//		IsLeaf: node.IsLeaf(),
	//		Radius: radius,
	//	})
	//})
	//
	//bytes, _ := json.Marshal(snapshots)
	//f, err := os.Create("C:\\Users\\kjk83317\\Desktop\\Unitiy\\Interpolation\\Assets\\Saves\\aabbs.json")
	//if err != nil {
	//	panic(err)
	//}
	//w := bufio.NewWriter(f)
	//w.WriteString(string(bytes))
	//w.Flush()
}
