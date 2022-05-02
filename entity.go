package main

import (
	"fmt"
	graphics "pokemon/Graphics"
	scripts "pokemon/Scripter"
	"time"

	"github.com/faiface/pixel"
)

type Entity struct {
	//Logic
	AttachedScript *scripts.Script

	//Graphics
	Sprite         *graphics.SpriteGroup
	frameToRender  string
	
	//Position things
	x, y           float64
	targetX,targetY float64


	//Clock things
	clockStart     time.Time
	clockTime      time.Duration
	clockActive    bool
}

func (e *Entity) Draw(t pixel.Target) {

	e.Sprite.Sprites[e.frameToRender].Draw(t, pixel.V(e.x, e.y), ImageScale)
}

func (e *Entity) Update(se *scripts.ScriptEngine) error {
	//Clock logic
	if e.clockActive {
		if time.Since(e.clockStart) > e.clockTime {
			e.AttachedScript.Resume()
			e.clockActive = false
		}
	}
	//Execute script
	err := se.ContinueScript(e.AttachedScript)
	
	speed:=1.0/60.0
	//Update position
	if e.targetX!=e.x{
		dir:=sign(e.targetX - e.x)
		e.x+=dir*speed
	}
	return err
}

func sign(n float64) float64{
	if n<0{
		return -1
	}
	return 1
}

func (e *Entity) SetPos(x, y float64) {
	e.AttachedScript.SetMemory(".tx", fmt.Sprint(x))
	e.AttachedScript.SetMemory(".ty", fmt.Sprint(x))
	e.targetX = x
	e.targetY = y
	e.x=x
	e.y=y	
}
