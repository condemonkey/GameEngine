package core

import (
	"game-engine/core/geo"
	"game-engine/math64"
	"game-engine/math64/vector3"
	"math/rand"
	"testing"
	"time"
)

func TestTree(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	min, max := -2000, 2000
	ecount := 5000

	tree := NewBVTree()

	var entities []Collider
	for i := 0; i < ecount; i++ {
		collider := NewCollider(geo.NewSphere(1)) // 30
		pointx := float64(math64.RandInt(min, max))
		pointy := float64(math64.RandInt(min, max))
		pointz := float64(math64.RandInt(min, max))
		collider.SetPosition(vector3.Vector3{X: pointx, Y: pointy, Z: pointz})
		collider.SetCollisionHandle(func(hit Collider) {
		})
		entities = append(entities, collider)
	}

	collider := NewCollider(geo.NewSphere(500)) // 30
	collider.SetPosition(vector3.Zero)

	entities = append(entities, collider)

	now := time.Now()
	for _, entity := range entities {
		tree.AddCollider(entity)
	}
	t.Log(time.Now().Sub(now).Milliseconds())

	loop := 20
	//wg := new(sync.WaitGroup)
	//for i := 0; i < loop; i++ {
	//	count := 0
	//	now = time.Now()
	//	for j := 0; j < 5000; j++ {
	//		wg.Add(1)
	//		go func() {
	//			count = tree.Intersect(collider)
	//			wg.Done()
	//		}()
	//	}
	//	wg.Wait()
	//	t.Log("detect time", time.Now().Sub(now).Milliseconds(), "collide", count)
	//	time.Sleep(time.Millisecond * 30)
	//}

	//for i := 0; i < loop; i++ {
	//	count := 0
	//	now = time.Now()
	//	for j := 0; j < 5000; j++ {
	//		count = tree.Intersect(collider)
	//	}
	//	t.Log("detect time2", time.Now().Sub(now).Milliseconds(), "collide", count)
	//	time.Sleep(time.Millisecond * 30)
	//}

	for i := 0; i < loop; i++ {
		count := 0
		now = time.Now()
		for j := 0; j < 5000; j++ {
			count = tree.Intersect2(collider)
		}
		t.Log("detect time3", time.Now().Sub(now).Milliseconds(), "collide", count)
		time.Sleep(time.Millisecond * 30)
	}

	for i := 0; i < loop; i++ {
		count := 0
		now = time.Now()
		for j := 0; j < 5000; j++ {
			count = tree.IntersectBp2(collider)
		}
		t.Log("detect time4", time.Now().Sub(now).Milliseconds(), "collide", count)
		time.Sleep(time.Millisecond * 30)
	}

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
	//		radius = node.Collider().(*SphereCollider).Radius()
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

	//objs := map[int]*GameObject{}
	//
	//now := time.Now()
	//for i := 0; i < 2000; i++ {
	//	obj := NewGameObject(
	//		NewTransform(vector3.Vector3{X: float64(i), Y: float64(i), Z: float64(i)}, vector3.One),
	//		//NewTransform(vector3.Zero, vector3.One),
	//		NewSphereCollider(vector3.Zero, 1))
	//	objs[i] = obj
	//}
	//t.Log(time.Now().Sub(now).Milliseconds())
	//
	//now = time.Now()
	//for _, obj := range objs {
	//	tree.AddCollider(obj.collider)
	//}
	//t.Log(time.Now().Sub(now).Milliseconds())

	//now = time.Now()
	//for i := 0; i < 2000; i++ {
	//	next := vector3.Add(objs[i].Transform().Position(), vector3.Vector3{2, 2, 2})
	//	objs[i].Transform().SetPosition(next)
	//}
	//t.Log(time.Now().Sub(now).Milliseconds())
	//
	//now = time.Now()
	//tree.Update()
	//t.Log(time.Now().Sub(now).Milliseconds())
}
