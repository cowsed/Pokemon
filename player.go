package main

import (
	"fmt"
	graphics "pokemon/Graphics"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/inkyblackness/imgui-go"
)

type Direction uint8

const (
	NoDirection Direction = iota
	Left
	Right
	Up
	Down
)

var DirNames = [5]string{"None", "Left", "Right", "Up", "Down"}
var speed = 3.0 / 60.0

type Player struct {
	spriteSheet *graphics.SpriteGroup
	spriteName  string

	x, y             float64
	targetX, targetY int

	currentDirection Direction
	queuedDirection  Direction
}

func (p *Player) DrawUI() {
	imgui.Text(fmt.Sprintf("pos: (%.2f, %.2f)", p.x, p.y))

	imgui.Text(fmt.Sprintf("posI: (%v, %v)", int(p.x), int(p.y)))

	_, _, tx, ty := p.CalcTargetPos()

	imgui.Text(fmt.Sprintf("target: (%v, %v)", tx, ty))

	dir := DirNames[p.currentDirection]
	qdir := DirNames[p.queuedDirection]
	imgui.Text(fmt.Sprintf("Going %s\tQueued: %s", dir, qdir))

	imgui.Text(fmt.Sprintf("Snap x: %.2f", abs(p.x-float64(tx))))
}

func (p *Player) CalcTargetPos() (float64, float64, int, int) {
	velX := 0.0
	velY := 0.0
	targetX := (p.x)
	targetY := (p.y)

	switch p.currentDirection {
	case Left:
		velX = -1
		targetX -= speed
	case Right:
		velX = 1
		targetX += speed * 2
	case Up:
		velY = 1
		targetY += speed * 2
	case Down:
		velY = -1
		targetY -= speed
	}

	return velX, velY, int(targetX), int(targetY)
}

func (p *Player) HandleInput(dir Direction) {
	//set if not moving, queue if already going somewhere
	if p.currentDirection == NoDirection {
		//if no direction currently we are snapped to the grid
		p.currentDirection = dir

		switch dir {
		case Left:
			p.targetX = int(p.x - 1)
		case Right:
			p.targetX = int(p.x + 1)
		case Up:
			p.targetY = int(p.y + 1)
		case Down:
			p.targetY = int(p.y - 1)
		}

	} else {
		p.queuedDirection = dir
	}

}

func (p *Player) Update() {

	snapTolerance := speed
	targetX := p.targetX
	targetY := p.targetY

	switch p.currentDirection {
	case Left:
		p.x -= speed
	case Right:
		p.x += speed

	case Up:
		p.y += speed
	case Down:
		p.y -= speed
	}

	if abs(p.x-float64(targetX)) < snapTolerance && abs(p.x-float64(targetX)) > 0 {
		//Snap to place
		p.x = float64(targetX)

		//queued and current are the same, dont do it
		if p.currentDirection == p.queuedDirection {
			p.queuedDirection = NoDirection
		}
		p.currentDirection = NoDirection
		//cycle through input buffer
		p.HandleInput(p.queuedDirection)

		p.queuedDirection = NoDirection
	}

	if abs(p.y-float64(targetY)) < snapTolerance && abs(p.y-float64(targetY)) > 0 {
		//Snap to place
		p.y = float64(targetY)

		//queued and current are the same, dont do it
		if p.currentDirection == p.queuedDirection {
			p.queuedDirection = NoDirection
		}
		p.currentDirection = NoDirection
		//cycle through input buffer
		p.HandleInput(p.queuedDirection)
		p.queuedDirection = NoDirection
	}

}

func (p *Player) Draw(win *pixelgl.Window) {

	m := abs((float64(p.targetY) - p.y) + (float64(p.targetX) - p.x))
	suffixIndex := int(m*4) % 4
	suffix := []string{"1", "2", "1", "3"}[suffixIndex]

	//Frame code shamelessly stolen from entity - NoDirection will return the direction from the previous frame
	dirName := []string{p.spriteName[:len(p.spriteName)-1], "left", "right", "up", "down"}[p.currentDirection]

	p.spriteName = dirName + suffix

	p.spriteSheet.Sprites[p.spriteName].DrawWorldPosition(win, pixel.ZV, ImageScale)
}
