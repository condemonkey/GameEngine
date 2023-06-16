package core

import (
	"game-engine/core/transform"
	"game-engine/math64/vector3"
)

type Transform struct {
	BaseComponent
	parent    *Transform
	transform *transform.Transform
}

func NewTransform() *Transform {
	return &Transform{
		transform: transform.NewTransform(),
	}
}

func (t *Transform) Awake() {
}

func (t *Transform) Scale() vector3.Vector3 {
	return t.transform.Scale
}

func (t *Transform) SetScale(scale vector3.Vector3) {
	t.transform.Scale = scale
}

func (t *Transform) SetParent(parent *Transform) {
	t.parent = parent
}

func (t *Transform) InternalTransform() *transform.Transform {
	return t.transform
}

func (t *Transform) SetPosition(pos vector3.Vector3) {
	t.transform.Position = pos
	if t.Entity().Active() {
		t.entity.SetAttribute(IsTransformDirty, true)
	}
}

func (t *Transform) SetRotation(rot vector3.Vector3) {
	t.transform.Rotation = rot
}

func (t *Transform) Position() vector3.Vector3 {
	return t.transform.Position
}

func (t *Transform) DeltaDistance(target *Transform) float64 {
	return 0
}
