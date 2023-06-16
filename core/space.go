package core

import (
	"game-engine/core/collision"
	"sync"
)

type Space struct {
	tree *collision.BVTree
	wg   *sync.WaitGroup
}

func NewSpace() *Space {
	space := &Space{}
	space.wg = new(sync.WaitGroup)
	space.tree = collision.NewBVTree(space.wg)
	return space
}
