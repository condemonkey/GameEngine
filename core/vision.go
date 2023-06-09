package core

import (
	"game-engine/core/collision"
)

// 네트워크 플레이어들만 가지고 있으면 될 듯
// 새로운 entity가 등장하면 me(target추가 모든 object)<->target(add player 해당 player들한테만 자신의 변경 사항을 broadcasting 하면 됨)
type Vision struct {
	BaseComponent
	UpdatableComponent
	radius          float64
	visibleEntities []collision.Collider
	Collider        *Collider
	tree            *collision.BVTree
}

func NewVision(radius float64, tree *collision.BVTree) *Vision {
	return &Vision{
		radius:   radius,
		Collider: NewShapeCollider(NewSphere(radius)),
		tree:     tree,
	}
}

func (t *Vision) Awake() {
	t.Collider.SetTransform(t.Transform())
	//t.Entity().Collider().AddCollisionHandler(t)
}

func (t *Vision) Update(dt int) {
	// 범위 안 object들을 검출
	t.visibleEntities = t.tree.Query(t.Collider)
}

func (t *Vision) VisibleEntities() []collision.Collider {
	return t.visibleEntities
}

//CollisionEnter(collider *Collider)
//CollisionExit(collider *Collider)
