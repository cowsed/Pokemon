package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Scene struct {
	Env      *TiledEnvironment
	Entities map[string]Entity
}

type CollisionData int

const (
	NotPassable CollisionData = iota
	Level1
	Level2
	Level3
	Level4
	Level5
	Level6
	HeightChange
	Water
)

var CollisionAlpha float64 = 1

var NotImplemented = pixel.RGBA{R: 1, G: 0, B: 1, A: CollisionAlpha}

var CollisionColors = []pixel.RGBA{
	{R: .6, G: 0, B: 0, A: CollisionAlpha},
	{R: 0, G: .6, B: 0, A: CollisionAlpha},
	NotImplemented,
	NotImplemented,
	NotImplemented,
	NotImplemented,
	NotImplemented,
	NotImplemented,
	NotImplemented,
	{R: 1, G: 0, B: 1, A: CollisionAlpha},
}

var M = [][]CollisionData{
	{0, 0, 0, 0, 0, 0, 1, 8, 8, 8, 8, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 8, 8, 8, 8, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 8, 8, 8, 8, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
	{0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
	{0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
}

type TiledEnvironment struct {
}

func (t *TiledEnvironment) GetCollision(x, y int) CollisionData {
	if x < 0 || y < 0 {
		return Level1
	} else if x >= len(M[0]) || y >= len(M) {
		return NotPassable
	}
	return M[y][x]
}

func (t *TiledEnvironment) DrawCollisions(win *pixelgl.Window, campos pixel.Vec) {

	imd := imdraw.New(nil)
	imd.Clear()

	for x := 0; x < 24; x++ {
		for y := 0; y < 20; y++ {
			open := Game.TileOpen(x, y, 4, 4, "")

			//tile := NotPassable //t.GetCollision(x, y)
			col := pixel.RGBA{
				R: 1,
				G: 0,
				B: 0,
				A: CollisionAlpha,
			} ///CollisionColors[int(tile)]

			if open {
				//tile = Level1
				col = pixel.RGBA{
					R: 0,
					G: 1,
					B: 0,
					A: CollisionAlpha,
				}
			}

			p1 := pixel.V(float64(x), float64(y)).Scaled(16 * ImageScale).Add(campos)
			p2 := pixel.V(float64(x+1), float64(y+1)).Scaled(16 * ImageScale).Add(campos)

			imd.SetColorMask(color.RGBA{
				R: 120,
				G: 120,
				B: 120,
				A: 120,
			})
			imd.Color = col
			imd.Push(p1, p2)
			//imd.Circle(10, 10)
			imd.Rectangle(0)
		}
	}

	imd.Draw(win)
	imd.Clear()

}
