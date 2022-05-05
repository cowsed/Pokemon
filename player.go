package main

import (
	graphics "pokemon/Graphics"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Direction uint8

const (
	NoDirection Direction = iota
	Left
	Right
	Up
	Down
)

type Player struct {
	spriteSheet     *graphics.SpriteGroup
	x, y            float64
	queuedDirection Direction
}

func (p *Player) Update() {
	velX := 0.0
	velY := 0.0
	switch p.queuedDirection {
	case Left:
		velX = -1
	case Right:
		velX = 1
	case Up:
		velY = 1
	case Down:
		velY = -1
	}

	snapTolerance := 1.0 / 16.0
	speed := .8 / 60.0

	p.x += velX * speed
	p.y += velY * speed

	targetX := int(p.x) + int(velX)

	if abs(p.x-float64(targetX)) < snapTolerance {
		p.x = float64(targetX)
		p.queuedDirection = NoDirection
	}

	//targetY := int(p.y) + int(velY)
	//if abs(p.y)-abs(float64(targetY)) < snapTolerance {
	//	p.y = float64(targetY)
	//	p.queuedDirection = NoDirection
	//}

}

func (p *Player) Draw(win *pixelgl.Window) {
	middleX := win.Bounds().W() / 2
	middleY := win.Bounds().H() / 2

	//x := middleX/(16*ImageScale) - .5
	//y := middleY/(16*ImageScale) - .5

	p.spriteSheet.Sprites["up1"].DrawScreenPosition(win, pixel.V(middleX, middleY+1.25*16), ImageScale)
	//p.spriteSheet.Sprites["up1"].DrawWorldPosition(win, pixel)
}
