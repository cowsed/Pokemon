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

func CheckBlocked(currentTile, nextTile CollisionData) bool {

	if currentTile == nextTile {
		return true
	}
	if nextTile == NotPassable {
		return false
	}
	if nextTile == Water {
		return false
	}

	return true
}

const NoKey = pixelgl.Key0

var FramesToHold = 8

var DirectionCounts map[pixelgl.Button]int = make(map[pixelgl.Button]int)

func directionFromKey(k pixelgl.Button) Direction {
	var directionToGive Direction
	switch k {
	case pixelgl.KeyUp:
		directionToGive = Up
	case pixelgl.KeyDown:
		directionToGive = Down
	case pixelgl.KeyLeft:
		directionToGive = Left
	case pixelgl.KeyRight:
		directionToGive = Right

	}
	return directionToGive
}

func (p *Player) HandleAllInput() {
	toWatch := []pixelgl.Button{pixelgl.KeyUp, pixelgl.KeyDown, pixelgl.KeyLeft, pixelgl.KeyRight}
	for _, k := range toWatch {
		if Game.ui.Pressed(k) {
			DirectionCounts[k]++
		} else {
			delete(DirectionCounts, k)
		}

		directionToGive := directionFromKey(k)
		if DirectionCounts[k] > FramesToHold {
			p.faceDir(directionToGive)
			p.HandleHeldInput(directionToGive)

		} else if DirectionCounts[k] > 0 {
			p.faceDir(directionToGive)
		}
	}
}

func (p *Player) faceDir(d Direction) {
	switch d {
	case Up:
		p.faceUp()
	case Down:
		p.faceDown()
	case Left:
		p.faceLeft()
	case Right:
		p.faceRight()
	}
}
func (p *Player) faceDown() {
	p.spriteName = "down1"
}
func (p *Player) faceUp() {
	p.spriteName = "up1"
}
func (p *Player) faceLeft() {
	p.spriteName = "left1"
}
func (p *Player) faceRight() {
	p.spriteName = "right1"
}
func (p *Player) HandleHeldInput(dir Direction) {

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

		if !CheckBlocked(Game.CurrentScene.Env.GetCollision(int(p.x), int(p.y)), Game.CurrentScene.Env.GetCollision(p.targetX, p.targetY)) {
			p.targetX = int(p.x)
			p.targetY = int(p.y)
			p.currentDirection = NoDirection
		}

	} else {
		if dir != p.currentDirection {
			p.queuedDirection = dir
		}
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
		p.HandleHeldInput(p.queuedDirection)

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
		p.HandleHeldInput(p.queuedDirection)
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

func (p *Player) DrawUI() {
	imgui.Text(fmt.Sprintf("pos: (%.2f, %.2f)", p.x, p.y))

	imgui.Text(fmt.Sprintf("posI: (%v, %v)", int(p.x), int(p.y)))

	_, _, tx, ty := p.CalcTargetPos()

	imgui.Text(fmt.Sprintf("target: (%v, %v)", tx, ty))

	dirName := DirNames[p.currentDirection]
	qdirName := DirNames[p.queuedDirection]
	imgui.Text(fmt.Sprintf("Going %s\tQueued: %s", dirName, qdirName))

	imgui.Text(fmt.Sprintf("Snap x: %.2f", abs(p.x-float64(tx))))

	collisionTile := Game.CurrentScene.Env.GetCollision(int(p.x), int(p.y))
	imgui.Text(fmt.Sprint("Current Collision Data:", []string{"Not Passable", "Level 1", "Level 2", "Level 3", "Level 4", "Level 5", "Level 6", "Water"}[collisionTile]))

	toWatch := []pixelgl.Button{pixelgl.KeyUp, pixelgl.KeyDown, pixelgl.KeyLeft, pixelgl.KeyRight}
	imgui.Separator()
	imgui.Text("Keys")
	for _, k := range toWatch {
		imgui.Text(fmt.Sprintf("%v - %d", k, DirectionCounts[k]))
	}

}
