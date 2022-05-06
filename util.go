package main

import "github.com/faiface/pixel/pixelgl"

func (g *GameStruct) UpdateSize() {

	w := g.win.Bounds().Max.X - g.win.Bounds().Min.X
	h := g.win.Bounds().Max.Y - g.win.Bounds().Min.Y
	g.WordHandler.Drawer.UpdateSize(w, h)
}

func (g *GameStruct) CheckWindowUpdates() {
	wasResized := CheckIfResized(g.win)
	if wasResized {
		Game.UpdateSize()
	}

}

//Assorted helper stuff that doesn't really involve the game
var lastW float64
var lastH float64

func CheckIfResized(win *pixelgl.Window) bool {
	w := win.Bounds().Max.X - win.Bounds().Min.X
	h := win.Bounds().Max.Y - win.Bounds().Min.Y
	if w != lastW || h != lastH {
		lastW = w
		lastH = h
		return true
	}
	return false
}
