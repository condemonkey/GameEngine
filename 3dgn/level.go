package main

import (
	"game-engine/core"
)

type D3GnLevel struct {
	*core.BaseLevel
}

func NewLevel() *D3GnLevel {
	return &D3GnLevel{
		BaseLevel: core.NewBaseLevel(),
	}
}

func (b *D3GnLevel) Load() {

}

func (b *D3GnLevel) Run() {

}
