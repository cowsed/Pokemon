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

	ScriptEngine  *scripts.ScriptEngine
	ActiveScripts map[string]*scripts.Script
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
	g.ScriptEngine = scripts.NewDefaultScriptEngine()
	//Register all the custom functions
	g.ScriptEngine.RegisteredFunction("dblog", DebugLogFunction)
	g.ScriptEngine.RegisteredFunction("dblogf", DebugLogFFunction)
	g.ScriptEngine.RegisteredFunction("say", DialogueFunction)
	g.ScriptEngine.RegisteredFunction("sayf", DialogueFFunction)

	//Debug program
	scr1 := scripts.NewScriptFromFile("Resources/Scripts/factorial.ps")
	scr1.Resume()
	g.ActiveScripts["db"] = scr1
}
