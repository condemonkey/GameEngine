package test

import (
	"game-engine/core"
	"game-engine/math64"
	"game-engine/math64/vector3"
	"math/rand"
	"time"
)

type Level1 struct {
	*core.BaseLevel
}

func (l *Level1) Load() {
	rand.Seed(time.Now().UnixNano())

	min, max := -2000, 2000
	ymin, ymax := 10, 10
	ecount := 7000

	for i := 0; i < ecount; i++ {
		entity := l.CreateEmptyEntity("ghost", "enemy")

		entity.AddComponent(core.NewSphereCollider(1)) // 크기1
		entity.AddComponent(core.NewVision())
		entity.AddComponent(core.NewTransform())
		entity.AddComponent(NewGhost())

		x := float64(math64.RandInt(min, max))
		y := float64(math64.RandInt(ymin, ymax))
		z := float64(math64.RandInt(min, max))
		entity.Transform().SetPosition(vector3.Vector3{X: x, Y: y, Z: z})

		l.AddEntity(entity)
	}

	entity := l.CreateEmptyEntity("player", "player")
	entity.AddComponent(core.NewSphereCollider(1)) // 크기1
	entity.AddComponent(core.NewVision())
	entity.AddComponent(core.NewTransform())
	entity.AddComponent(core.NewNetworkSession()) // 네트워크 통신을 위한 컴포넌트 해당 컴포넌트로 Send/Read?

	l.AddEntity(entity)
}

func (l *Level1) Create() core.Level {
	return &Level1{
		BaseLevel: core.NewBaseLevel("TestLevel"),
	}
}

func (l *Level1) OnLoad() {
}

func (l *Level1) OnDestroy() {
}
