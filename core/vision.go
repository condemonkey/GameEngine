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
	collision.Hittable
	visibles map[int]*NetworkSession
}

func NewVision() *Vision {
	return &Vision{
		visibles: make(map[int]*NetworkSession),
	}
}

func (t *Vision) AddVisible(session *NetworkSession) {
	t.visibles[session.Entity().Id()] = session
}

func (t *Vision) RemoveVisible(session *NetworkSession) {
	delete(t.visibles, session.Entity().Id())
}

func (t *Vision) Position() vector3.Vector3 {
	return t.Transform().Position()
}

func (t *Vision) Scale() vector3.Vector3 {
	return t.Transform().Scale()
}

// 시야 범위 안 network 전송
func (t *Vision) SendVisible() {
	fmt.Println("SendVisible NetworkSession", len(t.visibles))
	for _, net := range t.visibles {
		net.Send()
	}
}

//func (t *Vision) OnHit(handle collision.Hittable) {
//	entity := handle.(*Vision).Entity()
//	fmt.Println(entity.Transform().Position())
//}
