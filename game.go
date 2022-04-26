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
	WordHandler *DialogueHandler
	atlas       *text.Atlas
	win         *pixelgl.Window
	ui          *pixelui.UI
	logger      *Logger

	ScriptEngine  *scripts.ScriptHandler
	ActiveScripts map[string]*scripts.Script
}

func (g *GameStruct) UpdateSize() {

	w := g.win.Bounds().Max.X - g.win.Bounds().Min.X
	h := g.win.Bounds().Max.Y - g.win.Bounds().Min.Y
	g.WordHandler.drawer.UpdateSize(w, h)
}

func (g *GameStruct) HandleInput() {
	if g.ui.JustPressed(pixelgl.MouseButtonLeft) {
		fmt.Fprintln(g.logger, "Clicky")
	}
	if g.ui.JustPressed(pixelgl.KeyEnter) {
		g.WordHandler.HandleKey(pixelgl.KeyEnter)
	}

}
func (g *GameStruct) Draw(win *pixelgl.Window) {
	win.Clear(color.Black)
	g.WordHandler.Draw(win)
	g.ui.Draw(win)
}

func (g *GameStruct) InitializeScriptEngine() {
	g.ScriptEngine = scripts.NewDefualtScriptEngine()
	//Register all the custom functions
	g.ScriptEngine.RegisteredFunction("dblog", DebugLogFunction)
	g.ScriptEngine.RegisteredFunction("dblogf", DebugLogFFunction)
	g.ScriptEngine.RegisteredFunction("say", DialogueFunction)
	g.ScriptEngine.RegisteredFunction("sayf", DialogueFFunction)

	//fmt.Printf("g.ScriptEngine.RegisteredFunctions: %v\n", g.ScriptEngine.RegisteredFunctions)

	g.ScriptEngine.PrintFunc = func(s string) { fmt.Fprintln(Game.WordHandler.drawer.text, s) }

	//Debug program
	scr1 := scripts.NewScript(DebugScriptSource)
	scr1.Resume()
	g.ActiveScripts["db"] = scr1
}

func (g *GameStruct) CheckWindowUpdates() {
	wasResized := CheckIfResized(g.win)
	if wasResized {
		Game.UpdateSize()
	}

}
