package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	scripter "pokemon/Scripter"
)

var Game GameStruct

func run() {
	//Setup Window
	Game = GameStruct{
		WordHandler:   nil,
		ActiveScripts: map[string]*scripter.Script{},
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

	//Game loop
	for !win.Closed() {

		//Handle Input
		Game.HandleInput()

		//Execute game logic in scripts
		err = Game.ScriptEngine.ContinueScript(Game.ActiveScripts["db"])

		if err != nil {
			fmt.Fprintf(Game.logger, "Error executing script: %v", err.Error())
		}

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
