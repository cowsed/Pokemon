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

//say
var DialogueFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: 1,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {
		txt := args[0]
		Game.WordHandler.SetText(txt, script)
		script.Pause()
		return nil
	},
	Docstring: "say string	:	outputs string to the dialogue box",
}

//sayf
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
	Docstring: "stops the script for time specified as the arguement (literal or variable). Do not use this to wait to get to a location or something. It will have adverse effects. Use it for narrative timing moments. (scripted look in circle",
}

var MoveFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: 2,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {

		return nil
	},
	Docstring: "",
}
var SetPosFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: 2,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptEngine) error {

		return nil
	},
	Docstring: "",
}
