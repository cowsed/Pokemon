package main

import (
	graphics "pokemon/Graphics"
	scripts "pokemon/Scripter"

	"github.com/faiface/pixel"
)

type Entity struct {
	AttachedScript *scripts.Script
	Sprite         *graphics.SpriteGroup
	frameToRender  string
	x, y           float64
}

func (e *Entity) Draw(t pixel.Target) {
	e.Sprite.Sprites[e.frameToRender].Draw(t, pixel.V(e.x, e.y), ImageScale)
}
