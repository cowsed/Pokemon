package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
)

type TextDrawer struct {
	orig            pixel.Vec //top left
	bounds          pixel.Vec //bottom-right
	backgroundColor pixel.RGBA
	bg              *imdraw.IMDraw
	text            *text.Text
}

func NewTextDrawer(atlas *text.Atlas, windowW, windowH float64) *TextDrawer {
	basicTxt := text.New(pixel.V(0, 0), atlas)

	td := TextDrawer{
		orig:            pixel.Vec{X: 0, Y: 0},
		bounds:          pixel.Vec{X: 400, Y: 0},
		backgroundColor: pixel.RGB(.95, .8, .8),
		bg:              imdraw.New(nil),
		text:            basicTxt,
	}
	td.text.Color = pixel.RGB(1, 0, 0)
	td.MakeTextBox()
	td.UpdateSize(windowW, windowH)
	return &td
}
func (td *TextDrawer) UpdateSize(windowW, windowH float64) {
	borderY := 10.0
	borderX := 20.0
	h := windowH / 3
	td.orig = pixel.V(borderX, h-borderY)
	td.bounds = pixel.V(windowW-borderX, borderY)
	td.MakeTextBox()

}
func (td *TextDrawer) MakeTextBox() {
	td.bg = imdraw.New(nil)

	td.bg.Color = td.backgroundColor
	td.bg.Push(td.orig)

	td.bg.Color = td.backgroundColor
	td.bg.Push(pixel.V(td.orig.X, td.bounds.Y))

	td.bg.Color = td.backgroundColor
	td.bg.Push(td.bounds)

	td.bg.Color = td.backgroundColor
	td.bg.Push(pixel.V(td.bounds.X, td.orig.Y))
	td.bg.Polygon(0)

}

func (td *TextDrawer) Draw(t pixel.Target) {
	td.bg.Draw(t)

	textBL := td.orig.Add(pixel.V(20,-100))
	td.text.Draw(t, pixel.IM.Moved(textBL).Scaled(textBL, 2))
}
