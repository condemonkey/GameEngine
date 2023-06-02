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

	var entities []*Entity

	entitiesCount := 1000
	for i := 0; i < entitiesCount; i++ {
		entity := NewEntity(nil, "test", "")
		entity.AddComponent(NewSphereCollider())
		entity.AddComponent(NewTransform())

		pointx := float64(math64.RandInt(0, 300))
		pointy := float64(math64.RandInt(0, 300))
		pointz := float64(math64.RandInt(0, 300))
		entity.Transform().SetPosition(vector3.Vector3{X: pointx, Y: pointy, Z: pointz})

		entity.awake()
		entities = append(entities, entity)
	}

	tree := collision.NewBVTree()
	now := time.Now()
	for _, entity := range entities {
		tree.AddCollider(entity.Collider())
	}
	t.Log(time.Now().Sub(now).Milliseconds())

	//now = time.Now()
	//for _, entity := range entities {
	//	point := float64(math64.RandInt(0, 20))
	//	entity.Transform().SetPosition(vector3.Sumf(entity.Transform().Position(), point))
	//	tree.UpdateCollider(entity.Collider())
	//}
	//t.Log(time.Now().Sub(now).Milliseconds())

	snapshots := collision.Snapshots{}

	tree.Traverse(func(node *collision.BVTreeNode) {
		aabb := node.FatAABB()
		radius := float64(0)
		if node.IsLeaf() {
			radius = node.Collider().(*Collider).InternalShape().(*Sphere).Radius()
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
