package main

import (
	"fmt"
	"log"
	graphics "pokemon/Graphics"
	scripts "pokemon/Scripter"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var Game GameStruct
var logger *Logger

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

	logger = &Logger{
		internal: "",
	}
	log.SetOutput(logger)
	log.Println("Beginning log")

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
		//Handle Input
		Game.HandleInput()

		//Do Scripts
		for _, name := range getActiveEntityNames(&Game) {
			err = Game.ActiveEntites[name].DoScript(Game.ScriptEngine)

			if err != nil {
				log.Printf("Error executing script of entity %s: %v\n", name, err.Error())
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
