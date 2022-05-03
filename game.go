package main

import (
	"image/color"
	"log"

	scripts "pokemon/Scripter"

	"github.com/dusk125/pixelui"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

var TestEntity *Entity

const ImageScale float64 = 5

type GameStruct struct {
	WordHandler *DialogueHandler
	atlas       *text.Atlas
	win         *pixelgl.Window
	ui          *pixelui.UI
	//logger      *Logger

	env Environment

	ScriptEngine  *scripts.ScriptEngine
	ActiveEntites map[string]*Entity
}

func (g *GameStruct) HandleInput() {
	if g.ui.JustPressed(pixelgl.MouseButtonLeft) {
		v := g.win.MousePosition()
		log.Println("Clicky", v)
	}
	if g.ui.JustPressed(pixelgl.KeyEnter) {
		g.WordHandler.HandleKey(pixelgl.KeyEnter)
	}
}

func (g *GameStruct) Draw(win *pixelgl.Window) {
	win.Clear(color.Gray{
		Y: 80,
	})
	//Draw environment
	g.env.Draw(win, pixel.V(16*ImageScale, 0))

	//Draw entities
	for _, name := range getActiveEntityNames(g) {
		g.ActiveEntites[name].Draw(win)

	}

	//Draw box over active entity
	if selectedEntity != "" {
		e := Game.ActiveEntites[selectedEntity]
		d := imdraw.New(nil)

		d.Color = color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 60,
		}
		//d.SetMatrix(pixel.IM.Moved(pixel.V(e.x, e.y)))
		d.Push(pixel.V((e.x+.5)*16*ImageScale, (e.y+.5)*16*ImageScale))
		d.Circle(70, 10)

		d.Draw(win)
	}
	//Draw game ui
	g.WordHandler.Draw(win)

	//Draw debug ui
	g.ui.Draw(win)
}

func (g *GameStruct) InitializeGraphics() {

	var err error
	g.env, err = NewImageEnvFromFile("/home/rich/SelfGaming/Pokemon/Resources/Environments/PalleteTown/ptown.png")
	if err != nil {
		panic(err)
	}
}

func (g *GameStruct) InitializeScriptEngine() {
	g.ScriptEngine = scripts.NewDefaultScriptEngine()
	//Register all the custom functions
	g.ScriptEngine.RegisterFunction("dblog", DebugLogFunction)
	g.ScriptEngine.RegisterFunction("dblogf", DebugLogFFunction)
	g.ScriptEngine.RegisterFunction("say", DialogueFunction)
	g.ScriptEngine.RegisterFunction("sayf", DialogueFFunction)
	g.ScriptEngine.RegisterFunction("clearmem", ClearMemFunction)
	g.ScriptEngine.RegisterFunction("setframe", SetFrameFunction)
	g.ScriptEngine.RegisterFunction("wait", WaitFunction)
	g.ScriptEngine.RegisterFunction("setpos", SetPosFunction)
	g.ScriptEngine.RegisterFunction("movx", MovXFunction)
	g.ScriptEngine.RegisterFunction("movy", MovYFunction)

}
func (g *GameStruct) AddEntity(name string, E *Entity) {
	g.ActiveEntites[name] = E
	E.AttachedScript.SetMemory(".name", name)
}
