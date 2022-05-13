package main

import (
	"fmt"
	"math"
	graphics "pokemon/Graphics"
	scripts "pokemon/Scripter"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var EntityWalkSpeed = 2.0 / 60.0

type Entity struct {
	name string
	//Logic
	AttachedScript *scripts.Script

	//Graphics
	Sprite        *graphics.SpriteGroup
	frameToRender string

	//Position things
	moving           bool
	x, y             float64
	targetX, targetY float64

	//Clock things
	clockStart  time.Time
	clockTime   time.Duration
	clockActive bool

	hidden bool
}

func NewEntity(name string, x, y float64, initialFrame, sheetPath, sheetDataPath string, script *scripts.Script) (*Entity, error) {
	sprite, err := graphics.LoadSprite(sheetPath, sheetDataPath)
	if err != nil {
		return nil, err
	}
	script.Resume()
	e := Entity{
		name:           name,
		AttachedScript: script,
		Sprite:         sprite,
		frameToRender:  initialFrame,
		x:              x,
		y:              y,
		targetX:        x,
		targetY:        y,
	}
	return &e, nil
}

func (e *Entity) Interact() {
	if e.AttachedScript.HasLabel("onInteract") {

		e.AttachedScript.Call("onInteract")

	} else {
		fmt.Println("No interac")
	}
}

func (e *Entity) Draw(t *pixelgl.Window, offset pixel.Vec) {
	if e.hidden {
		return
	}
	e.Sprite.Sprites[e.frameToRender].DrawWorldPosition(t, pixel.V(e.x, e.y).Add(offset), ImageScale)
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
	if err != nil {
		fmt.Println(err)
	}

	//Update position X
	e.HandleMovement()

	return err
}
func (e *Entity) HandleMovement() {

	//X
	deltaX := e.targetX - e.x
	if Game.TileOpen(int(e.targetX), int(e.targetY), int(e.x), int(e.y), e.name) {
		if abs(deltaX) > 1.0/16 { //If off by more than a pixel

			dir := sign(deltaX)
			e.x += dir * EntityWalkSpeed

			//Calculate sprite
			//left foot - center - right foot - center
			_, r := math.Modf(abs(e.x))
			frameNum := int(r * 4)
			directionName := []string{"left", "right"}[int(dir+1)/2]
			suffix := []string{"2", "1", "3", "1"}[frameNum]

			e.frameToRender = directionName + suffix
			e.moving = true

		} else if deltaX != 0 { //Close enough to finish
			//Back to normal position
			dir := sign(deltaX)
			directionName := []string{"left", "right"}[int(dir+1)/2]
			e.frameToRender = directionName + "1"

			e.moving = false
			//Snap to pixel perfect location
			e.x = e.targetX
			e.AttachedScript.Resume()
		}
	}

	//Y
	deltaY := e.targetY - e.y
	if Game.TileOpen(int(e.targetX), int(e.targetY), int(e.x), int(e.y), e.name) {
		if abs(deltaY) > 1.0/16 { //If off by more than a pixel
			dir := sign(deltaY)
			e.y += dir * EntityWalkSpeed

			//Calculate sprite
			//left foot - center - right foot - center
			_, r := math.Modf(abs(e.y))
			frameNum := int(r * 4)
			directionName := []string{"down", "up"}[int(dir+1)/2]
			suffix := []string{"2", "1", "3", "1"}[frameNum]

			e.frameToRender = directionName + suffix
			e.moving = true

		} else if deltaY != 0 { //Close enough to finish
			//Back to normal position
			dir := sign(deltaY)
			directionName := []string{"down", "up"}[int(dir+1)/2]
			e.frameToRender = directionName + "1"

			//Snap to pixel perfect location
			e.y = e.targetY
			e.moving = false
			e.AttachedScript.Resume()
		}
	}
}

func abs(n float64) float64 {
	if n < 0 {
		return -n
	}
	return n
}
func sign(n float64) float64 {
	if n < 0 {
		return -1
	}
	return 1
}

func (e *Entity) SetPos(x, y float64) {
	e.targetX = x
	e.targetY = y
	e.x = x
	e.y = y
}
