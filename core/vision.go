package core

import (
	"fmt"
	"game-engine/core/collision"
	"game-engine/math64/vector3"
)

// 네트워크 플레이어들만 가지고 있으면 될 듯
// 새로운 entity가 등장하면 me(target추가 모든 object)<->target(add player 해당 player들한테만 자신의 변경 사항을 broadcasting 하면 됨)
type Vision struct {
	BaseComponent
	UpdatableComponent
	collision.Hittable
	visibleEntities []collision.Collider
	query           collision.TreeQuery
	radius          float64
	cnt             int
}

func NewVision(radius float64) *Vision {
	return &Vision{
		radius: radius,
	}
}

func (t *Vision) Awake() {
	t.query = t.Level().TreeQuery()
}

func (t *Vision) Update(dt int) {
	//t.Transform().SetPosition(vector3.Zero)
}

func (t *Vision) FinalUpdate(dt int) {
	//t.query.IntersectRangeAsync(t.Transform().Position(), t.radius, func(hits int) {
	//
	//})
	t.query.IntersectRangeCollidersAsync(t.Transform().Position(), t.radius, func(hits []*collision.Collider) {
	})
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
