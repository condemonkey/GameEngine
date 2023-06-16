package core

import (
	"game-engine/math64"
	"game-engine/math64/vector3"
	"math/rand"
	"testing"
	"time"
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

var defaultLevel = &Level1{
	BaseLevel: NewBaseLevel("Test"),
}

func TestLevel(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	min, max := -2000, 2000
	ymin, ymax := 10, 10
	ecount := 8000

	for i := 0; i < ecount; i++ {
		entity := defaultLevel.CreateEmptyEntity("test", "player")
		entity.AddComponent(NewSphereCollider(1)) // 크기1
		entity.AddComponent(NewVision())          // 시야 테스트
		entity.AddComponent(NewTransform())

		x := float64(math64.RandInt(min, max))
		y := float64(math64.RandInt(ymin, ymax))
		z := float64(math64.RandInt(min, max))
		entity.Transform().SetPosition(vector3.Vector3{X: x, Y: y, Z: z})

		defaultLevel.AddEntity(entity)
	}

	end := time.Now()
	for i := 0; i < 10000; i++ {
		start := time.Now()
		defaultLevel.update(int(start.Sub(end).Milliseconds()))
		end = time.Now()

		t.Log("update delay", end.Sub(start).Milliseconds())

		time.Sleep(time.Millisecond * 100) // 0.1초 주기로 갱신
	}

	defaultLevel.SaveSnapshot()
}
