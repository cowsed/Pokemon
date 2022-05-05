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

type Player struct {
	spriteSheet      *graphics.SpriteGroup
	x, y             float64
	currentDirection Direction
	queuedDirection  Direction
}

func (p *Player) DrawUI() {
	imgui.Text(fmt.Sprintf("pos: (%.2f,%.2f)", p.x, p.y))
}

func (p *Player) Update() {
	velX := 0.0
	velY := 0.0
	switch p.currentDirection {
	case Left:
		velX = -1
	case Right:
		velX = 1
	case Up:
		velY = 1
	case Down:
		velY = -1
	}

	speed := 5 / 60.0
	snapTolerance := speed

	p.x += velX * speed
	p.y += velY * speed

	targetX := int(p.x) + int(velX)

	if abs(p.x-float64(targetX)) < snapTolerance {
		//Snap to place
		p.x = float64(targetX)

		//queued and current are the same, dont do it
		if p.currentDirection == p.queuedDirection {
			p.queuedDirection = NoDirection
		}
		p.currentDirection = NoDirection
		//cycle through input buffer
		p.currentDirection = p.queuedDirection
		p.queuedDirection = NoDirection
	}

	//targetY := int(p.y) + int(velY)
	//if abs(p.y-float64(targetY)) < snapTolerance {
	//	//Snap to place
	//	p.y = float64(targetY)
	//	p.currentDirection = NoDirection
	//	//cycle through input buffer
	//	p.currentDirection = p.queuedDirection
	//	p.queuedDirection = NoDirection
	//}
	//targetY := int(p.y) + int(velY)
	//if abs(p.y)-abs(float64(targetY)) < snapTolerance {
	//	p.y = float64(targetY)
	//	p.queuedDirection = NoDirection
	//}

}

func (p *Player) Draw(win *pixelgl.Window) {
	middleX := win.Bounds().W() / 2
	middleY := win.Bounds().H() / 2

	p.spriteSheet.Sprites["up1"].DrawScreenPosition(win, pixel.V(middleX, middleY+1.25*16), ImageScale)
}
