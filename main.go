package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	scripter "pokemon/Scripter"
)

/*set $x 0
set $y 1
set $z 0
dblog "wow - from a script"
set $counter 12
dblogf 2 "Printing %s nums" $counter
startlabel:
addI $counter -1 $counter
addI $x $y $z
set $x $y
set $y $z
dblogf 2 %s $x
jmpne startlabel $counter 0*/
var DebugScriptSource = `
say "wow, text"
say "wow, text. but again"

`

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
		Bounds:    pixel.R(0, 0, 1800, 1000),
		Resizable: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	check(err)
	Game.win = win

	//Initialize the game engine
	Game.InitializeUI()
	Game.InitializeGameUI()
	Game.InitializeScriptEngine()

	//Game loop
	for !win.Closed() {

		//Handle Input
		Game.HandleInput()

		//Execute game logic in scripts
		err = Game.ScriptEngine.ContinueScript(Game.ActiveScripts["db"])

		if err != nil {
			fmt.Fprintf(Game.logger, "Error executing script: %s", err.Error())
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
