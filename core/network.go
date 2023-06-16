package core

import (
	"game-engine/core/collision"
)

type NetworkSession struct {
	BaseComponent
	UpdatableComponent
	visionRange float64
	query       collision.TreeQuery
	visibles    []*NetworkSession
	//grpc eventstream
	//grpc owner
}

func NewNetworkSession() *NetworkSession {
	return &NetworkSession{
		BaseComponent: NewBaseComponent(),
	}
}

func (t *NetworkSession) Awake() {
	t.query = t.Level().SpaceQuery()
}

func (t *NetworkSession) Update(dt int) {
}

func (t *NetworkSession) FinalUpdate(dt int) {
	// vision update
	t.query.CollidingBox(t.Transform().Position(), t.visionRange, t.UpdateVisible)
}

func (t *NetworkSession) UpdateVisible(colliders []collision.Collider) {
	for _, col := range colliders {
		entity := col.(*Collider).Entity()
		if entity.NetworkSession() != nil {
			// 내 시야에 들어온 NetworkSession 추가, 각 Session마다 진행
			t.Vision().AddVisible(entity.NetworkSession())
		} else {
			// 상대방 시야에 나를 추가하자
			entity.Vision().AddVisible(t)
			// 몬스터의 경우 이동 -> visible.SendVisible() 하면 동기화 해야 할 플레이어 들한테만 날아간다.
		}
	}
}

// 네트워크 전송
func (t *NetworkSession) Send() {
}
