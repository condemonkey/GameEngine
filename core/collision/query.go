package collision

import "game-engine/math64/vector3"

type TreeQuery interface {
	CollidingSphereAsync(center vector3.Vector3, distance float64, result SearchCallback)
	CollidingSphere(origin vector3.Vector3, distance float64, result SearchCallback)
	// 가장 성능이 빠름
	CollidingBox(origin vector3.Vector3, size float64, result SearchCallback)
	// 레이캐스팅
	RayCasting(ray Ray, result SearchCallback)
}
