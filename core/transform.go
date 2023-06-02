package core

import (
	"game-engine/math64/vector3"
)

type Transform struct {
	BaseComponent
	position vector3.Vector3
	rotation vector3.Vector3
	scale    vector3.Vector3
	parent   *Transform
}

func NewTransform() *Transform {
	return &Transform{
		scale:    vector3.One,
		position: vector3.Zero,
		rotation: vector3.Zero,
	}
}

func (t *Transform) Awake() {
}

func (t *Transform) Scale() vector3.Vector3 {
	return t.scale
}

func (t *Transform) SetScale() vector3.Vector3 {
	return t.scale
}

func (t *Transform) SetParent(parent *Transform) {
	t.parent = parent
}

func (t *Transform) SetPosition(pos vector3.Vector3) {
	//t.entity.SetAttribute(IsTransformDirty, true)
	t.position = pos
}

func (t *Transform) SetRotation(rot vector3.Vector3) {
}

func (t *Transform) Position() vector3.Vector3 {
	return t.position
}

func (t *Transform) DeltaDistance(target *Transform) float64 {
	return 0
}
