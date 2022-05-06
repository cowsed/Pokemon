package main

import (
	"fmt"
	"log"
	"os"
	graphics "pokemon/Graphics"
	scripts "pokemon/Scripter"
	"runtime/pprof"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var Game GameStruct
var logger *Logger
var pc *PerformanceCounter = &PerformanceCounter{
	lastAddTime:   time.Now(),
	lastFPS:       -99,
	lastFrameTime: -99,
	doingVsync:    true,
}

var UpdatesPerFrame = 1

func run() {
	//Setup Window
	Game = GameStruct{
		WordHandler:   nil,
		ActiveEntites: map[string]*Entity{},
	}

	cfg := pixelgl.WindowConfig{
		Title:     "In Dev",
		Icon:      []pixel.Picture{},
		Bounds:    pixel.R(0, 0, 1800, 900),
		Resizable: true,
		VSync:     pc.doingVsync,
	}

	win, err := pixelgl.NewWindow(cfg)
	check(err)
	Game.win = win

	//Initialize the game engine
	Game.InitializeLogger()

	Game.InitializeGraphics()
	Game.InitializeGameUI()
	Game.InitializeUI()

	Game.InitializeScriptEngine()

	Game.LoadPlayer()

	//Debug environment
	{

		Game.env, err = NewImageEnvFromFile("Resources/Environments/PalleteTown/ptown.png")
		check(err)
	}

	//Debug program
	{
		for y := 0; y < 2*60; y++ {
			for x := 0; x < 2*80; x++ {
				TestEntity = &Entity{
					AttachedScript: nil,
					Sprite:         nil,
					frameToRender:  "down1",
				}
				TestEntity.x = float64(x)
				TestEntity.y = float64(y)
				TestEntity.targetX = float64(x)
				TestEntity.targetY = float64(y)

				scr1 := scripts.NewScriptFromFile("Resources/Scripts/math.ps")
				scr1.Resume()
				TestEntity.AttachedScript = scr1

				TestEntity.Sprite, err = graphics.LoadSprite("Resources/Sprites/Builtin/officer.png", "Resources/Sprites/Builtin/officer.json")
				check(err)

				Game.AddEntity(fmt.Sprintf("point(%v, %v)", x, y), TestEntity)
			}
		}
	}

	//Game loop
	for !win.Closed() {
		//Handle Input
		for i := 0; i < UpdatesPerFrame; i++ {
			Game.HandleInput()

			//Do Scripts
			//for _, name := range getActiveEntityNames(&Game) {
			for name, e := range Game.ActiveEntites {

				err = e.Update(Game.ScriptEngine)

				if err != nil {
					log.Printf("Error executing script of entity '%s': %v\n", name, err.Error())
				}
			}
			Game.player.Update()
		}
		//Draw
		Game.Draw(win)

		//Draw debug ui
		Game.DrawDebugUI()
		Game.ui.Draw(win)

		//Checks for window resizing and such
		Game.CheckWindowUpdates()
		win.Update()

		pc.DoCount()
	}

}
func main() {
	f, _ := os.Create("prof.pprof")
	pprof.StartCPUProfile(f)
	defer func() {
		pprof.StopCPUProfile()
		f.Close()
	}()

	pixelgl.Run(run)

	fmt.Println("Exitting")
	SaveNote()
	fmt.Println("Finished Shutdown")

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

/*
	//Boy
	TestEntity = &Entity{
		AttachedScript: nil,
		Sprite:         nil,
		frameToRender:  "down1",
	}
	TestEntity.x = 9
	TestEntity.y = 3
	TestEntity.targetX = 9
	TestEntity.targetY = 3

	scr1 := scripts.NewScriptFromFile("Resources/Scripts/spin.ps")
	scr1.Resume()
	TestEntity.AttachedScript = scr1

	TestEntity.Sprite, err = graphics.LoadSprite("Resources/Sprites/Builtin/brendan.png", "Resources/Sprites/Builtin/brendan.json")
	check(err)

	Game.AddEntity("test_guy", TestEntity)

	//Officer
	var TestEntity2 = &Entity{
		AttachedScript: nil,
		Sprite:         nil,
		frameToRender:  "down2",
	}
	TestEntity.x = 3
	TestEntity.y = 3
	TestEntity.targetX = 3
	TestEntity.targetY = 3

	scr2 := scripts.NewScriptFromFile("Resources/Scripts/math.ps")
	TestEntity2.AttachedScript = scr2
	scr2.Resume()

	TestEntity2.Sprite, err = graphics.LoadSprite("Resources/Sprites/Builtin/officer.png", "Resources/Sprites/Builtin/officer.json")
	check(err)

	Game.AddEntity("test_guy2", TestEntity2)
*/
