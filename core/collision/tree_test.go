package collision

import (
	"bufio"
	"encoding/json"
	"game-engine/core/geo"
	"game-engine/math64"
	"game-engine/math64/vector3"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestCollision(t *testing.T) {
	tree := NewBVTree()
	collider1 := NewCollider(geo.NewSphere(0.5))
	collider1.SetPosition(vector3.Vector3{X: 0, Y: 0, Z: 0})
	collider1.SetCollisionHandle(func(hit Collider) {
	})

	tree.AddCollider(collider1)

	collider := NewCollider(geo.NewSphere(0.5))
	collider.SetPosition(vector3.Vector3{X: 0.6, Y: 0.6, Z: 0.6})
	collider.SetCollisionHandle(func(hit Collider) {
	})

	//abs := tree.Query(collider)
	//for _, col := range abs {
	//	col.Intersect(collider)
	//}
	//t.Log(abs)
	//
	//result := collider.Intersect(collider1)
	//t.Log(result)
}

func TestUpdatePosition(t *testing.T) {
	tree := NewBVTree()
	collider1 := NewCollider(geo.NewSphere(0.5))
	collider1.SetPosition(vector3.Vector3{X: 0, Y: 0, Z: 0})
	collider1.SetCollisionHandle(func(hit Collider) {
	})

	tree.AddCollider(collider1)

	collider := NewCollider(geo.NewSphere(0.5))
	collider.SetPosition(vector3.Vector3{X: 0.6, Y: 0.6, Z: 0.6})
	collider.SetCollisionHandle(func(hit Collider) {
	})

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
	ecount := 5000

	tree := NewBVTree()

	var entities []Collider
	for i := 0; i < ecount; i++ {
		collider := NewCollider(geo.NewSphere(0.5)) // 30
		pointx := float64(math64.RandInt(min, max))
		pointy := float64(math64.RandInt(min, max))
		pointz := float64(math64.RandInt(min, max))
		collider.SetPosition(vector3.Vector3{X: pointx, Y: pointy, Z: pointz})
		collider.SetCollisionHandle(func(hit Collider) {
		})
		entities = append(entities, collider)
	}

	collider := NewCollider(geo.NewSphere(200)) // 30
	collider.SetPosition(vector3.Zero)

	entities = append(entities, collider)

	now := time.Now()
	for _, entity := range entities {
		tree.AddCollider(entity)
	}
	t.Log(time.Now().Sub(now).Milliseconds())

	var iters []*Iterator
	for i := 0; i < 5000; i++ {
		iters = append(iters, tree.NewIterator())
	}

	//wg := new(sync.WaitGroup)
	for i := 0; i < 30; i++ {
		now = time.Now()
		for j := 0; j < 5000; j++ {
			<-iters[j].SearchAsync(collider)
			//wg.Add(1)
			//index := j
			//go func() {
			//	iters[index].Search(collider)
			//	wg.Done()
			//}()
		}
		//wg.Wait()
		t.Log("detect time", time.Now().Sub(now).Milliseconds())
		//t.Log("collision count", count)
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

	snapshots := Snapshots{}

	tree.Traverse(func(node *BVTreeNode) {
		aabb := node.FatAABB()
		radius := float64(0)
		if node.IsLeaf() {
			radius = node.Collider().(*SphereCollider).Radius()
		}
		snapshots.Snapshots = append(snapshots.Snapshots, Snapshot{
			Min:    aabb.Min(),
			Max:    aabb.Max(),
			Center: aabb.Center(),
			Size:   aabb.Size(),
			IsLeaf: node.IsLeaf(),
			Radius: radius,
		})
	})

	bytes, _ := json.Marshal(snapshots)
	f, err := os.Create("C:\\Users\\kjk83317\\Desktop\\Unitiy\\Interpolation\\Assets\\Saves\\aabbs.json")
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	w.WriteString(string(bytes))
	w.Flush()

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
