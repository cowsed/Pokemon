package main

import (
	"fmt"
	"log"
	scripts "pokemon/Scripter"
	"strconv"
	"strings"
	"time"
)

/*
Create the interface through which the scripting language can communicate with the game
*/
//Print to the debug console
var DebugLogFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: 1,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {
		log.Println(args[0])
		return nil
	},
	Docstring: "Logs the specified text to the game console",
}

//Printf but to the debug console
var DebugLogFFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: -1,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {
		fstring := args[0]
		values := make([]interface{}, len(args)-1)
		for i := 0; i < len(values); i++ {
			values[i] = script.ParseValue(args[i+1])
		}

		log.Printf(fstring+"\n", values...)
		return nil
	}, Docstring: "Logs the specified text to the game console using printf. See the Sayf documentation",
}

//say msg
var DialogueFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: 1,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {
		txt := args[0]
		Game.WordHandler.SetText(txt, script)
		Game.lastInputHandle = Game.InputHandleFunc
		Game.InputHandleFunc = Game.WordHandler.HandleAllInput
		Game.WordHandler.OnClose = func() {
			Game.InputHandleFunc = Game.lastInputHandle
		}
		script.Pause()
		return nil
	},
	Docstring: "say string	:	outputs string to the dialogue box",
}

//sayf n fstring args
var DialogueFFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: -1,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {
		//Extract format specifier and values
		fstring := args[0]
		values := make([]interface{}, len(args)-1)
		for i := 0; i < len(values); i++ {
			values[i] = script.ParseValue(args[i+1])
		}
		txt := fmt.Sprintf(fstring, values...)

		//Update text
		Game.WordHandler.SetText(txt, script)
		Game.lastInputHandle = Game.InputHandleFunc
		Game.InputHandleFunc = Game.WordHandler.HandleAllInput
		Game.WordHandler.OnClose = func() {
			Game.InputHandleFunc = Game.lastInputHandle
		}
		script.Pause()
		return nil
	},
	Docstring: "sayf n format ...	:	outputs formatted string to the dialogue box. n is the number of variables + 1. i.e. sayf 2 %s $variable",
}

//clearmem
var ClearMemFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: 0,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {
		//Find all special variables (starting with .)
		specialStore := map[string]string{}
		for k, v := range script.Memory() {
			if strings.HasPrefix(k, ".") {
				specialStore[k] = v
			}
		}
		script.ClearMemory()
		for k, v := range specialStore {
			script.SetMemory(k, v)
		}
		return nil
	},
	Docstring: "Resets the scripts memory. Use this if your script relies on empty memory",
}

//setframe who frame
var SetFrameFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: 2,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {
		who := script.ParseValue(args[0])
		towhat := script.ParseValue(args[1])
		if _, ok := Game.ActiveEntites[who].Sprite.Sprites[towhat]; ok {
			Game.ActiveEntites[who].frameToRender = towhat
		}
		return nil
	},
	Docstring: "Sets the entity specified by argument 1 to the frame specified by argument 2. If the specified frame or entity does not exist, it does nothing",
}

//wait t
var WaitFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: 1,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {
		arg := script.ParseValue(args[0])
		t, _ := strconv.ParseFloat(arg, 64)
		script.Pause()
		tm := time.Duration(t * float64(time.Second))

		Game.ActiveEntites[script.GetMemory(".name")].clockStart = time.Now()
		Game.ActiveEntites[script.GetMemory(".name")].clockTime = tm
		Game.ActiveEntites[script.GetMemory(".name")].clockActive = true
		return nil
	},
	Docstring: "stops the script for time specified as the arguement (literal or variable).  Use it for narrative timing moments. (scripted look in circle)",
}

//setpos who x y
var SetPosFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: 3,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {
		name := script.ParseValue(args[0])
		xStr := args[1]
		yStr := args[2]
		x, err := strconv.ParseFloat(xStr, 64)
		if err != nil {
			return err
		}
		y, err := strconv.ParseFloat(yStr, 64)
		if err != nil {
			return err
		}

		Game.ActiveEntites[name].SetPos(x, y)

		return nil
	},
	Docstring: "",
}

//movx who tx
var MovXFunction = scripts.ScriptFunction{
	NumArguments: 2,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {

		who := script.ParseValue(args[0])
		newXStr := args[1]
		newX, err := strconv.ParseFloat(newXStr, 64)

		if err != nil {
			return err
		}
		if _, ok := Game.ActiveEntites[who]; !ok {
			return fmt.Errorf("no entity named %s", who)
		}
		e := Game.ActiveEntites[who]

		//Already there
		if Game.ActiveEntites[who].targetX == newX {
			return nil
		}
		//Set Direction
		if newX < Game.ActiveEntites[who].x {
			Game.ActiveEntites[who].frameToRender = "left1"
		} else {
			Game.ActiveEntites[who].frameToRender = "right1"
		}

		//If theres a player in the way, go back a step so that it is still interactable
		if !Game.TileOpen(int(newX), int(e.y), int(e.x), int(e.y), "") {
			script.Backup(3)
			return scripts.ErrorYield
		}

		Game.ActiveEntites[who].targetX = newX
		Game.ActiveEntites[who].AttachedScript.Pause()

		return nil
	},

	Docstring: "movx who target : moves the sprite specified to the location in the x direction",
}

//movy who ty
var MovYFunction = scripts.ScriptFunction{
	NumArguments: 2,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {

		who := script.ParseValue(args[0])
		newYStr := args[1]
		newY, err := strconv.ParseFloat(newYStr, 64)

		if err != nil {
			return err
		}
		if _, ok := Game.ActiveEntites[who]; !ok {
			return fmt.Errorf("no entity named %s", who)
		}
		e := Game.ActiveEntites[who]

		//Already there
		if e.targetY == newY {
			return nil
		}

		//Set Direction
		if newY < e.y {
			e.frameToRender = "down1"
		} else {
			e.frameToRender = "up1"
		}

		//If theres a player in the way, go back a step so that it is still interactable
		if !Game.TileOpen(int(e.x), int(newY), int(e.x), int(e.y), "") {
			script.Backup(3)
			return scripts.ErrorYield
		}

		Game.ActiveEntites[who].targetY = newY
		Game.ActiveEntites[who].AttachedScript.Pause()

		return nil
	},

	Docstring: "movy who target : moves the sprite specified to the location in the y direction",
}

//movx who tx
var GetPosFunction = scripts.ScriptFunction{
	NumArguments: 3,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {

		who := script.ParseValue(args[0])
		xLoc := args[1]
		yLoc := args[2]
		script.SetMemory(xLoc[1:], fmt.Sprint(Game.ActiveEntites[who].x))
		script.SetMemory(yLoc[1:], fmt.Sprint(Game.ActiveEntites[who].y))

		return nil
	},

	Docstring: "get who x y : stores the x and y values of the entity who in to memory locations called x and y",
}

//hide who
var HideFunction = scripts.ScriptFunction{
	NumArguments: 1,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {

		who := script.ParseValue(args[0])

		if _, ok := Game.ActiveEntites[who]; !ok {
			return fmt.Errorf("no entity named %s", who)
		}

		Game.ActiveEntites[who].hidden = true

		return nil
	},

	Docstring: "hide who : hides the sprite specified ",
}

//hide who
var ShowFunction = scripts.ScriptFunction{
	NumArguments: 1,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {

		who := script.ParseValue(args[0])

		if _, ok := Game.ActiveEntites[who]; !ok {
			return fmt.Errorf("no entity named %s", who)
		}

		Game.ActiveEntites[who].hidden = false

		return nil
	},

	Docstring: "hide who : hides the sprite specified ",
}
