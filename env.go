package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type Environment interface {
	Draw(win pixel.Target, at pixel.Vec)
}

type GridEnv struct {
	imd *imdraw.IMDraw
}

func (g *GridEnv) Draw(win pixel.Target, at pixel.Vec) {

	g.imd.Draw(win)
}

func (g *GridEnv) MakeGrid() {
	grid := imdraw.New(nil)

	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			if (x%2 == 0) != (y%2 == 0) {
				grid.Color = colornames.Red
			} else {
				grid.Color = colornames.Orange
			}
			grid.Push(pixel.V(float64(x*16)*ImageScale, float64(y*16)*ImageScale))
			grid.Push(pixel.V(float64((x+1)*16)*ImageScale, float64(y*16)*ImageScale))
			grid.Push(pixel.V(float64((x+1)*16)*ImageScale, float64((y+1)*16)*ImageScale))
			grid.Push(pixel.V(float64(x*16)*ImageScale, float64((y+1)*16)*ImageScale))
			grid.Polygon(0)
		}
	}
	g.imd = grid
}
