package transform

import "game-engine/math64/vector3"

type Transform struct {
	Position vector3.Vector3
	Rotation vector3.Vector3
	Scale    vector3.Vector3
}

func NewTransform() *Transform {
	return &Transform{
		Scale: vector3.One,
	}
}
