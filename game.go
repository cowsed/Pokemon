package main

import (
	"image/color"
	"log"

	graphics "pokemon/Graphics"
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

	env Environment

	player *Player

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

	if g.ui.Pressed(pixelgl.KeyLeft) {
		g.player.HandleInput(Left)
	}
	if g.ui.Pressed(pixelgl.KeyRight) {
		g.player.HandleInput(Right)
	}
	if g.ui.Pressed(pixelgl.KeyDown) {
		g.player.HandleInput(Down)
	}
	if g.ui.Pressed(pixelgl.KeyUp) {
		g.player.HandleInput(Up)
	}

}

func (g *GameStruct) Draw(win *pixelgl.Window) {
	win.Clear(color.Gray{
		Y: 80,
	})

	playerPos := pixel.V(g.player.x, g.player.y)

	cameraPosScreenSpace := playerPos.Scaled(-ImageScale * 16)
	//Draw environment
	g.env.Draw(win, cameraPosScreenSpace)

	//Draw entities
	for _, name := range getActiveEntityNames(g) {
		g.ActiveEntites[name].Draw(win, playerPos.Scaled(-1))

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
		d.Push(pixel.V((e.x+.5)*16*ImageScale, (e.y+.5)*16*ImageScale).Sub(playerPos.Scaled(16 * ImageScale)))
		d.Circle(70, 10)

		d.Draw(win)
	}
	g.player.Draw(win)

	//Draw game ui
	g.WordHandler.Draw(win)

	//Draw debug ui
	g.ui.Draw(win)
}

func (g *GameStruct) LoadPlayer() {
	ss, err := graphics.LoadSprite("Resources/Sprites/Builtin/brendan.png", "Resources/Sprites/Builtin/brendan.json")
	check(err)
	g.player = &Player{
		spriteSheet: ss,
		x:           0,
		y:           0,
	}
}

func (g *GameStruct) InitializeGraphics() {

	var err error
	g.env, err = NewImageEnvFromFile("/home/rich/SelfGaming/Pokemon/Resources/Environments/Grid/grid.png")
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

func (g *GameStruct) InitializeLogger() {
	//Init log
	logger = &Logger{
		internal: "",
	}
	log.SetOutput(logger)
	log.Println("Beginning log")
}
