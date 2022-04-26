package main

import (
	"fmt"
	"image/color"

	scripts "pokemon/Scripter"

	"github.com/dusk125/pixelui"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type GameStruct struct {
	DialogueHandler *TextDrawer
	Scripts         map[string]*scripts.Script
	atlas           *text.Atlas
	win             *pixelgl.Window
	ui              *pixelui.UI
	logger          *Logger
	ScriptEngine    *scripts.ScriptHandler
}

func (g *GameStruct) UpdateSize() {

	w := g.win.Bounds().Max.X - g.win.Bounds().Min.X
	h := g.win.Bounds().Max.Y - g.win.Bounds().Min.Y
	g.DialogueHandler.UpdateSize(w, h)
}

func (g *GameStruct) HandleInput() {
	if g.ui.JustPressed(pixelgl.MouseButtonLeft) {
		fmt.Fprintln(g.logger, "Clicky")
	}
}
func (g *GameStruct) Draw(win *pixelgl.Window) {
	win.Clear(color.Black)
	g.DialogueHandler.Draw(win)
	g.ui.Draw(win)
}

func (g *GameStruct) InitializeScriptEngine() {
	g.ScriptEngine = scripts.NewDefualtScriptEngine()
	//Register all the custom functions
	g.ScriptEngine.RegisteredFunction("dblog", DebugLogFunction)
	g.ScriptEngine.RegisteredFunction("dblogf", DebugLogFFunction)

	//fmt.Printf("g.ScriptEngine.RegisteredFunctions: %v\n", g.ScriptEngine.RegisteredFunctions)

	g.ScriptEngine.PrintFunc = func(s string) { fmt.Fprintln(Game.DialogueHandler.text, s) }

	//Debug program
	scr1 := scripts.NewScript(DebugScriptSource)
	scr1.Resume()
	g.Scripts["db"] = scr1
}

func (g *GameStruct) CheckWindowUpdates() {
	wasResized := CheckIfResized(g.win)
	if wasResized {
		Game.UpdateSize()
	}

}
