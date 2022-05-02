package main

import (
	graphics "pokemon/Graphics"
	scripts "pokemon/Scripter"
	"time"

	"github.com/faiface/pixel"
)

type Entity struct {
	AttachedScript *scripts.Script
	Sprite         *graphics.SpriteGroup
	frameToRender  string
	x, y           float64
	clockStart     time.Time
	clockTime      time.Duration
	clockActive    bool
}

func (e *Entity) Draw(t pixel.Target) {

	e.Sprite.Sprites[e.frameToRender].Draw(t, pixel.V(e.x, e.y), ImageScale)
}

func (e *Entity) DoScript(se *scripts.ScriptEngine) error {
	if e.clockActive {
		if time.Since(e.clockStart) > e.clockTime {
			e.AttachedScript.Resume()
			e.clockActive = false
		}
	}
	err := se.ContinueScript(e.AttachedScript)
	return err
}
