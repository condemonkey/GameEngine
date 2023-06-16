package test

import (
	"game-engine/core"
)

type Ghost struct {
	core.BaseComponent
	core.UpdatableComponent
}

func NewGhost() *Ghost {
	return &Ghost{
		BaseComponent: core.NewBaseComponent(),
	}
}

func (g *Ghost) Awake() {

}

func (g *Ghost) Update(dt int) {
	// 이동 일단 시키자
	position := g.Transform().Position()
	position.X += 0.01
	position.Y += 0.01
	position.Z += 0.01
	g.Transform().SetPosition(position)

	// 시야 범위 안 network session에 통지 하자.
	// Vision은 NetworkSession에서 동기화 시킨다.
	//g.Vision().SendVisible()
}

func (g *Ghost) FinalUpdate(dt int) {
	//g.Level().SpaceQuery().CollidingSphere(g.Transform().Position(), 50, func(handles []collision.Collider) {
	//
	//})
}
