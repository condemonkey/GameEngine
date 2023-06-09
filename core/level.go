package core

import (
	"game-engine/core/collision"
	"game-engine/core/geo"
)

type Level interface {
	NewEntity(name, tag string)
	AddEntity(entity *Entity)
	RemoveEntity(entity *Entity)
	FindEntitiesOverlapCollider(collider collision.Collider) []collision.Collider
	Load()
	OnLoad()
	OnDestroy()
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
	entity.AddComponent(NewCollider(collision.NewCollider(geo.NewSphere(0.5))))
	entity.AddComponent(NewTransform())
	return entity
}

func (b *BaseLevel) AddEntity(entity *Entity) {
	entity.awake()
	b.entities[entity.Id()] = entity
	b.tree.AddCollider(entity.Collider().CollisionCollider())
}

func (b *BaseLevel) RemoveEntity(entity *Entity) {
	delete(b.entities, entity.Id())
	b.tree.RemoveCollider(entity.Collider().CollisionCollider())
}

func (b *BaseLevel) FindEntitiesOverlapCollider(collider collision.Collider) []collision.Collider {
	return nil
	//return b.tree.Query(collider)
}

func (b *BaseLevel) Update(dt int) {
	b.entity.update(dt)
	for _, entity := range b.entities {
		//if entity.Attribute(IsTransformDirty) {
		//	// update transform position???????
		//	entity.SetAttribute(IsTransformDirty, false)
		//}
		entity.update(dt)
	}
}

func (b *BaseLevel) Load() {
}

func (b *BaseLevel) OnLoad() {
	b.entity.awake()
}

func (b *BaseLevel) OnDestroy() {
}
