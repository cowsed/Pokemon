package main

import (
	"fmt"
	scripts "pokemon/Scripter"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type DialogueHandler struct {
	WaitingForConfirmation bool
	ListedText             string
	Active                 bool
	ActiveScript           *scripts.Script
	drawer                 *TextDrawer
}

func (dh *DialogueHandler) HandleKey(button pixelgl.Button) {

	if dh.WaitingForConfirmation && button == pixelgl.KeyEnter {
		dh.Confirmed()
	}
}
func (dh *DialogueHandler) Confirmed() {
	dh.WaitingForConfirmation = false

	if dh.ActiveScript != nil {
		dh.ActiveScript.Resume()
	}
	dh.Close()

}
func (dh *DialogueHandler) Close() {
	dh.Active = false
	dh.ActiveScript = nil

}
func (dh *DialogueHandler) SetText(txt string, from *scripts.Script) {
	dh.ActiveScript = from

	dh.ListedText = txt
	dh.drawer.text.Clear()
	fmt.Fprint(dh.drawer.text, txt)

	dh.Active = true
	dh.WaitingForConfirmation = true
}
func (dh *DialogueHandler) Draw(win *pixelgl.Window) {
	if dh.Active {

		dh.drawer.Draw(win)
	}
}

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

	textBL := td.orig.Add(pixel.V(20, -100))
	td.text.Draw(t, pixel.IM.Moved(textBL).Scaled(textBL, 2))
}
