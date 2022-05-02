package main

import (
	"fmt"
	"image/color"

	scripts "pokemon/Scripter"

	"github.com/dusk125/pixelui"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

var grid *imdraw.IMDraw

var TestEntity *Entity

const ImageScale float64 = 5

type GameStruct struct {
	WordHandler *DialogueHandler
	atlas       *text.Atlas
	win         *pixelgl.Window
	ui          *pixelui.UI
	logger      *Logger

	ScriptEngine  *scripts.ScriptEngine
	ActiveEntites map[string]*Entity
}

func (g *GameStruct) HandleInput() {
	if g.ui.JustPressed(pixelgl.MouseButtonLeft) {
		v := g.win.MousePosition()
		fmt.Fprintln(g.logger, "Clicky", v)
	}
	if g.ui.JustPressed(pixelgl.KeyEnter) {
		g.WordHandler.HandleKey(pixelgl.KeyEnter)
	}
}

func (g *GameStruct) LoadGraphics() {
}

func (g *GameStruct) Draw(win *pixelgl.Window) {
	win.Clear(color.Gray{
		Y: 80,
	})
	grid.Draw(win)

	//sprite.Draw(win, pixel.IM.Scaled(pixel.V(0, 0), 1).Moved(pixel.V(8, 16)))
	for _, name := range getActiveEntityNames(g) {
		g.ActiveEntites[name].Draw(win)

	} //Officer.Sprites[FrameToRender].Draw(g.win, pixel.V(3, 2+1.0/16.0), ImageScale)

	g.WordHandler.Draw(win)

	g.ui.Draw(win)
}

func (g *GameStruct) InitializeGraphics() {
	g.LoadGraphics()
	g.MakeGrid()
}

func (g *GameStruct) InitializeScriptEngine() {
	g.ScriptEngine = scripts.NewDefaultScriptEngine()
	//Register all the custom functions
	g.ScriptEngine.RegisteredFunction("dblog", DebugLogFunction)
	g.ScriptEngine.RegisteredFunction("dblogf", DebugLogFFunction)
	g.ScriptEngine.RegisteredFunction("say", DialogueFunction)
	g.ScriptEngine.RegisteredFunction("sayf", DialogueFFunction)
	g.ScriptEngine.RegisteredFunction("clearmem", ClearMemFunction)
	g.ScriptEngine.RegisteredFunction("setframe", SetFrameFunction)
	g.ScriptEngine.RegisteredFunction("wait", WaitFunction)

}
func (g *GameStruct) AddEntity(name string, E *Entity) {
	g.ActiveEntites[name] = E
	E.AttachedScript.SetMemory(".name", name)
}

func (g *GameStruct) MakeGrid() {
	grid = imdraw.New(nil)
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
}
