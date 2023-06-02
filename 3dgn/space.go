package main

import (
	"game-engine/core/collision"
	g3ncore "github.com/g3n/engine/core"
)

type Space struct {
	tree  *collision.BVTree
	scene *g3ncore.Node
	//mesh  *graphic.Mesh
	//box   *graphic.Mesh
}
