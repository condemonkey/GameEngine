package main

import (
	"game-engine/core"
	g3ncore "github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

type GeoSphere struct {
	core.UpdatableComponent
	core.BaseComponent
	scene *g3ncore.Node
	mesh  *graphic.Mesh
	box   *graphic.Mesh
}

func NewGeoSphere(scene *g3ncore.Node) *GeoSphere {
	return &GeoSphere{
		scene: scene,
	}
}

func (t *GeoSphere) Awake() {
	collider := t.Entity().Collider()
	shape := collider.InternalShape().(*core.Sphere)
	{
		geom := geometry.NewSphere(shape.Radius(), 16, 8)
		mat := material.NewStandard(math32.NewColor("DarkBlue"))
		mat.SetWireframe(true)
		t.mesh = graphic.NewMesh(geom, mat)
		position := t.Transform().Position()
		t.mesh.SetPosition(float32(position.X), float32(position.Y), float32(position.Z))

		t.scene.Add(t.mesh)
	}
	{
		aabb := collider.AABB()
		geom := geometry.NewBox(float32(aabb.With()), float32(aabb.Height()), float32(aabb.Depth()))
		mat := material.NewStandard(math32.NewColor("DarkBlue"))
		mat.SetWireframe(true)
		t.box = graphic.NewMesh(geom, mat)
		position := t.Transform().Position()
		t.box.SetPosition(float32(position.X), float32(position.Y), float32(position.Z))
		t.scene.Add(t.box)
	}
}

func (t *GeoSphere) Update(dt int) {

}
