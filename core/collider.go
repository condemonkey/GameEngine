package core

import (
	"game-engine/core/collision"
)

type CollisionHandler interface {
	OnHit(entity *Entity)
}

type Collider struct {
	BaseComponent
	collision.Hittable
	collider  *collision.Collider
	listeners []CollisionHandler
}

func NewCollider(collider *collision.Collider) *Collider {
	return &Collider{
		BaseComponent: NewBaseComponent(),
		collider:      collider,
	}
}

func (p *Collider) Awake() {
	// component transform 연결
	p.collider.SetTransform(p.Transform().InternalTransform())
}

func (p *Collider) CollisionEnter(handle collision.Hittable) {
	for _, h := range p.listeners {
		h.OnHit(handle.(*Collider).Entity())
	}
}

func (p *Collider) CollisionCollider() *collision.Collider {
	return p.collider
}
