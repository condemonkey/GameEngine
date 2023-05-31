package core

type Collider struct {
	UpdatableComponent
	BaseComponent
}

func NewCollider() *Collider {
	return &Collider{
		BaseComponent: NewBaseComponent(),
	}
}

func NewSphereCollider() *Collider {
	return &Collider{
		BaseComponent: NewBaseComponent(),
	}
}

func (c *Collider) Update(dt int) {

}
