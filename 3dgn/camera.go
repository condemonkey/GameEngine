package main

import (
	"game-engine/core"
	"github.com/g3n/engine/camera"
)

type Camera struct {
	core.BaseComponent
}

func NewCamera() *Camera {
	return &Camera{}
}

func (t *Camera) Awake() {
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	//scene.Add(cam)
}
