package core

import (
	"game-engine/math64/vector3"
	"testing"
)

type Level1 struct {
	*BaseLevel
}

func (b *Level1) Load() {
}

func (b *Level1) OnLoad() {
}

func (b *Level1) OnDestroy() {
}

func TestLevel(t *testing.T) {
	level := &Level1{
		BaseLevel: NewBaseLevel(),
	}
	entity := level.CreateEntity("test", "test1level")
	entity.AddComponent(NewTransform())
	entity.AddComponent(NewSphereCollider())

	entity.Transform().SetPosition(vector3.One)

	level.AddEntity(entity)
}
