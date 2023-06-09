package geo

import "game-engine/math64/vector3"

type Box struct {
	Size vector3.Vector3
}

func (s Box) Type() ShapeType {
	return ShapeBox
}
