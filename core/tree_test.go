package core

import (
	"bufio"
	"encoding/json"
	"game-engine/core/collision"
	"game-engine/math64"
	"game-engine/math64/vector3"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestTree(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	tree := collision.NewBVTree()

	var entities []*Entity

	for i := 0; i < 2000; i++ {
		entity := NewEntity(nil, "test", "")
		entity.AddComponent(NewSphereCollider())
		entity.AddComponent(NewTransform())
		entity.AddComponent(NewVision(10, tree)) // 60

		pointx := float64(math64.RandInt(-1000, 1000))
		pointy := float64(math64.RandInt(-200, 200))
		pointz := float64(math64.RandInt(-1000, 1000))
		entity.Transform().SetPosition(vector3.Vector3{X: pointx, Y: pointy, Z: pointz})

		entities = append(entities, entity)
	}

	for i := 0; i < 30; i++ {
		entity := NewEntity(nil, "player", "")
		entity.AddComponent(NewSphereCollider())
		entity.AddComponent(NewTransform())
		entity.AddComponent(NewVision(50, tree))

		pointx := float64(math64.RandInt(-1000, 1000))
		pointy := float64(math64.RandInt(-200, 200))
		pointz := float64(math64.RandInt(-1000, 1000))
		entity.Transform().SetPosition(vector3.Vector3{X: pointx, Y: pointy, Z: pointz})

		entities = append(entities, entity)
	}

	now := time.Now()
	for _, entity := range entities {
		entity.awake()
		tree.AddCollider(entity.Collider())
		//comp := entity.FindComponent("Vision").(*Vision)
		//tree.AddCollider(comp.Collider)
	}
	t.Log(time.Now().Sub(now).Milliseconds())

	for i := 0; i < 30; i++ {
		now = time.Now()
		for _, entity := range entities {
			entity.Update(0)
		}
		t.Log(time.Now().Sub(now).Milliseconds())

		time.Sleep(time.Millisecond * 100)
	}

	for _, entity := range entities {
		comp := entity.FindComponent("Vision").(*Vision)
		tree.AddCollider(comp.Collider)
	}

	//col := NewEntity(nil, "test", "")
	//col.AddComponent(NewSphereCollider())
	//col.AddComponent(NewTransform())
	//pointx := float64(math64.RandInt(-50, 50))
	//pointy := float64(math64.RandInt(-50, 50))
	//pointz := float64(math64.RandInt(-50, 50))
	//col.Transform().SetPosition(vector3.Vector3{X: pointx, Y: pointy, Z: pointz})
	//col.Transform().SetScale(vector3.Vector3{X: 20, Y: 20, Z: 20})
	//col.awake()

	//now = time.Now()
	//collisions := tree.Query(col.Collider())
	//fmt.Println(len(collisions))
	//for _, shape := range collisions {
	//	//fmt.Println("on collide", shape.(*Collider).Shape().Center())
	//}
	//t.Log(time.Now().Sub(now).Milliseconds())

	//tree.AddCollider(col.Collider())
	//entities = append(entities, col)

	//now = time.Now()
	//for _, entity := range entities {
	//	point := float64(math64.RandInt(0, 20))
	//	entity.Transform().SetPosition(vector3.SumScalar(entity.Transform().Position(), point))
	//	tree.UpdateCollider(entity.Collider())
	//}
	//t.Log(time.Now().Sub(now).Milliseconds())

	snapshots := collision.Snapshots{}

	tree.Traverse(func(node *collision.BVTreeNode) {
		aabb := node.FatAABB()
		radius := float64(0)
		if node.IsLeaf() {
			radius = node.Collider().(*Collider).Shape().(*Sphere).Radius()
		}
		snapshots.Snapshots = append(snapshots.Snapshots, collision.Snapshot{
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
