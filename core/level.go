package core

import (
	"game-engine/core/collision"
)

type Level interface {
	NewEntity(name, tag string)
	AddEntity(entity *Entity)
	RemoveEntity(entity *Entity)
	RelocateEntity(entity *Entity)
	TreeQuery() collision.TreeQuery
	Load()
	OnLoad()
	OnDestroy()
	SaveSnapshot()
}

type BaseLevel struct {
	Level
	entity   *Entity
	entities map[int]*Entity
	tree     *collision.BVTree
}

func NewBaseLevel() *BaseLevel {
	b := &BaseLevel{}
	b.entity = NewEntity(b, "level", "")
	b.entities = make(map[int]*Entity)
	b.tree = collision.NewBVTree()
	return b
}

func (b *BaseLevel) AddLevelComponent(comp Component) {
	b.entity.AddComponent(NewTransform())
}

func (b *BaseLevel) CreateEntity(name, tag string) *Entity {
	entity := NewEntity(b, name, tag)
	entity.AddComponent(NewCollider(collision.NewSphereCollider(0.5))) // 디폴트.
	entity.AddComponent(NewTransform())
	return entity
}

func (b *BaseLevel) CreateEmptyEntity(name, tag string) *Entity {
	entity := NewEntity(b, name, tag)
	return entity
}

func (b *BaseLevel) AddEntity(entity *Entity) {
	// 모든 컴포넌트 초기화 및 활성화
	entity.awake()

	// 공간에 추가
	b.tree.AddCollider(entity.Collider().CollisionCollider())

	// 완료
	b.entities[entity.Id()] = entity
}

func (b *BaseLevel) RemoveEntity(entity *Entity) {
	delete(b.entities, entity.Id())
	b.tree.RemoveCollider(entity.Collider().CollisionCollider())
}

func (b *BaseLevel) RelocateEntity(entity *Entity) {
	b.tree.UpdateCollider(entity.Collider().CollisionCollider())
}

func (b *BaseLevel) TreeQuery() collision.TreeQuery {
	return b.tree
}

func (b *BaseLevel) Update(dt int) {
	//b.entity.update(dt)
	for _, entity := range b.entities {
		entity.update(dt)
		//if entity.Attribute(IsTransformDirty) {
		//	b.tree.UpdateCollider(entity.Collider().CollisionCollider())
		//	entity.SetAttribute(IsTransformDirty, false)
		//}
	}
	for _, entity := range b.entities {
		entity.finalUpdate(dt)
	}
	b.tree.WaitGroup()
}

func (b *BaseLevel) Load() {
}

func (b *BaseLevel) OnLoad() {
	b.entity.awake()
}

func (b *BaseLevel) OnDestroy() {
}

func (b *BaseLevel) SaveSnapshot() {
	b.tree.Snapshot()
}
