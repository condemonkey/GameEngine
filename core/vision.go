package core

import (
	"fmt"
	"game-engine/core/collision"
	"game-engine/core/geo"
	"game-engine/math64/vector3"
)

// 네트워크 플레이어들만 가지고 있으면 될 듯
// 새로운 entity가 등장하면 me(target추가 모든 object)<->target(add player 해당 player들한테만 자신의 변경 사항을 broadcasting 하면 됨)
type Vision struct {
	BaseComponent
	UpdatableComponent
	collision.Hittable
	visibleEntities []collision.Collider
	cube            collision.Collider
	tree            *collision.BVTree
}

func NewVision(radius float64, tree *collision.BVTree) *Vision {
	return &Vision{
		cube: collision.NewCollider(geo.Sphere{
			Radius: radius,
		}),
		tree: tree,
	}
}

func (t *Vision) Awake() {
	t.cube.SetHandle(t)
}

func (t *Vision) Update(dt int) {
	// 범위 안 object들을 검출
	t.tree.Intersect(t.cube)
	//t.visibleEntities = t.tree.Query(t.collider)
}

func (t *Vision) VisibleEntities() []collision.Collider {
	return t.visibleEntities
}

func (t *Vision) Position() vector3.Vector3 {
	return t.Transform().Position()
}

func (t *Vision) Scale() vector3.Vector3 {
	return t.Transform().Scale()
}

// on hit callback
func (t *Vision) OnHit(handle collision.Hittable) {
	entity := handle.(*Vision).Entity()
	fmt.Println(entity.Transform().Position())
}
