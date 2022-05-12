package main

import (
	"image/color"
	"log"

	graphics "pokemon/Graphics"
	scripts "pokemon/Scripter"
	ui "pokemon/UI"

	"github.com/dusk125/pixelui"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

var TestEntity *Entity

var ImageScale float64 = 5

type GameStruct struct {
	WordHandler *ui.DialogueHandler
	atlas       *text.Atlas
	win         *pixelgl.Window
	ui          *pixelui.UI

	env Environment

	player *Player

	ScriptEngine  *scripts.ScriptEngine
	ActiveEntites map[string]*Entity

	CurrentScene *Scene

	InputHandleFunc func()
	lastInputHandle func()
}

func (g *GameStruct) HandleInput() {
	if g.ui.JustPressed(pixelgl.MouseButtonLeft) {
		v := g.win.MousePosition()
		log.Println("Clicky", v)
	}

	g.InputHandleFunc()

}

func (g *GameStruct) Draw(win *pixelgl.Window) {
	win.Clear(pixel.RGBA{
		R: .3,
		G: .3,
		B: .3,
		A: 1,
	})
	win.SetComposeMethod(pixel.ComposeOver)

	playerPos := pixel.V(g.player.x, g.player.y)

	cameraPosScreenSpace := playerPos.Scaled(-ImageScale * 16)
	//Draw environment
	g.env.Draw(win, cameraPosScreenSpace)

	if showCollisionOverlay {
		g.CurrentScene.Env.DrawCollisions(win, cameraPosScreenSpace.Add(win.Bounds().Center()).Sub(pixel.V(7*ImageScale, 7*ImageScale)))
	}
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

		d.Push(pixel.V((e.x)*16*ImageScale, (e.y)*16*ImageScale).Sub(playerPos.Scaled(16 * ImageScale)).Add(win.Bounds().Center()))
		d.Circle(60, 10)

		d.Draw(win)
	}

	g.player.Draw(win)

	//Draw game ui
	g.WordHandler.Draw(win)

}

func (g *GameStruct) LoadPlayer() {
	ss, err := graphics.LoadSprite("Resources/Sprites/Builtin/may.png", "Resources/Sprites/Builtin/may.json")
	check(err)
	g.player = &Player{
		spriteSheet: ss,
		x:           4,
		y:           2,
		spriteName:  "down1",
	}
	g.InputHandleFunc = g.player.HandleAllInput
}

func (g *GameStruct) InitializeGraphics() {

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
	g.ScriptEngine.RegisterFunction("getpos", GetPosFunction)
	g.ScriptEngine.RegisterFunction("hide", HideFunction)
	g.ScriptEngine.RegisterFunction("show", ShowFunction)

}

//Entity Stuff
func (g *GameStruct) AddEntity(name string, E *Entity) {
	g.ActiveEntites[name] = E
	E.AttachedScript.SetMemory(".name", name)
}
func (g *GameStruct) InteractAt(x, y int) {
	for _, e := range g.ActiveEntites {
		if int(e.x+.5) == x && int(e.y+.5) == y {
			e.Interact()
		}
	}
}

func (g *GameStruct) InitializeLogger() {
	//Init log
	logger = &Logger{
		internal: "",
	}
	log.SetOutput(logger)
	log.Println("Beginning log")
}

type Logger struct {
	internal string
}

func (l *Logger) Write(bs []byte) (int, error) {
	l.internal += string(bs)
	return len(bs), nil
}
func (l *Logger) String() string {
	return l.internal
}
