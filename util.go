package main

import "github.com/faiface/pixel/pixelgl"

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
