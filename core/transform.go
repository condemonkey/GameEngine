package core

import "game-engine/core/math64/vector3"

type Transform struct {
	BaseComponent
	position vector3.Vector3
	rotation vector3.Vector3
	scale    vector3.Vector3 //부모 스케일이 있다면?
	parent   *Transform
}

func NewTransform() *Transform {
	return &Transform{}
}

func (t *Transform) Awake() {
}

func (t *Transform) Scale() vector3.Vector3 {
	return t.scale
}

func (t *Transform) SetPosition(pos vector3.Vector3) {
	t.entity.SetAttribute(IsTransformDirty, true)
	t.position = pos
}

func (t *Transform) SetRotation(rot vector3.Vector3) {
}
