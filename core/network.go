package core

import "game-engine/math64/vector3"

type NetworkTransform struct {
	BaseComponent
	transform *Transform
}

func NewNetworkTransform() *NetworkTransform {
	return &NetworkTransform{
		BaseComponent: NewBaseComponent(),
	}
}

func (n *NetworkTransform) Awake() {
	n.transform = n.Transform()
	//n.Vision().VisibleEntities()
	// 나한텐 일단 다 추가함
	// n.Entity().VisibleEntities().AddEntity(finder)
	// not is player -> 플레이어가 아니라면
	// entity.Vision().AddEntity(n.Entity()) 몬스터에 나를 추가 이동시 나한테 브로드캐스팅 오도록
}

func (n *NetworkTransform) SetPosition(pos vector3.Vector3) {
	n.transform.SetPosition(pos) // space tree rebalance
	//내 위치가 변경되면 시야 범위 내 entitie들을 갱신 한다.
	//entitiy.Vision().Refresh()?????
	//monster들은 player만 시야범위 안에 삽입하자.

	//시야에 속한 entity들 한테 broadcasting
	//n.Vision().VisibleEntities().Send()
}
