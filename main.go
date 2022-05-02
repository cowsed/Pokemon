package main

import (
	"fmt"
	graphics "pokemon/Graphics"
	scripts "pokemon/Scripter"
	"runtime"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var Game GameStruct

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
	}

	win, err := pixelgl.NewWindow(cfg)
	check(err)
	Game.win = win

	//Initialize the game engine
	Game.InitializeGraphics()
	Game.InitializeGameUI()
	Game.InitializeUI()

	Game.InitializeScriptEngine()

	//Debug program
	{
		//Boy
		TestEntity = &Entity{
			AttachedScript: nil,
			Sprite:         nil,
			frameToRender:  "down1",
			x:              4,
			y:              4,
		}

		scr1 := scripts.NewScriptFromFile("Resources/Scripts/animtest.ps")
		scr1.Resume()
		TestEntity.AttachedScript = scr1

		TestEntity.Sprite, err = graphics.LoadSprite("Resources/Sprites/Builtin/brendan.png", "Resources/Sprites/Builtin/brendan.json")
		check(err)

		Game.AddEntity("test_guy", TestEntity)

		//OFficer
		var TestEntity2 = &Entity{
			AttachedScript: nil,
			Sprite:         nil,
			frameToRender:  "down2",
			x:              2,
			y:              2,
		}

		scr2 := scripts.NewScriptFromFile("Resources/Scripts/spin.ps")
		scr2.Resume()
		TestEntity2.AttachedScript = scr2

		TestEntity2.Sprite, err = graphics.LoadSprite("Resources/Sprites/Builtin/officer.png", "Resources/Sprites/Builtin/officer.json")
		check(err)

		Game.AddEntity("test_guy2", TestEntity2)
	}

	//Game loop
	for !win.Closed() {
		fmt.Println(runtime.NumGoroutine())

		//Handle Input
		Game.HandleInput()

		//Do Scripts
		for _, name := range getActiveEntityNames(&Game) {

			err = Game.ScriptEngine.ContinueScript(Game.ActiveEntites[name].AttachedScript)
			if err != nil {
				fmt.Fprintf(Game.logger, "Error executing script of entity %s: %v", name, err.Error())
			}
		}

		//Draw
		Game.DrawDebugUI()
		Game.Draw(win)

		//Checks for window resizing and such
		Game.CheckWindowUpdates()
		win.Update()
	}

}
func main() {

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
