package main

import (
	"fmt"
	scripts "pokemon/Scripter"
)

/*
Create the interface through which the scripting language can communicate with the game
*/
//Print to the debug console
var DebugLogFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: 1,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptHandler) error {
		Game.logger.Write([]byte(args[0] + "\n"))
		return nil
	},
	Docstring: "Logs the specified text to the game console",
}

//Printf but to the debug console
var DebugLogFFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: -1,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptHandler) error {
		fstring := args[0]
		values := make([]interface{}, len(args)-1)
		for i := 0; i < len(values); i++ {
			values[i] = scr.ParseValue(args[i+1])
		}

		fmt.Fprintf(Game.logger, fstring+"\n", values...)
		return nil
	}, Docstring: "Logs the specified text to the game console using printf. See the Sayf documentation",
}

//say
var DialogueFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: 1,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptHandler) error {
		txt := args[0]
		Game.WordHandler.SetText(txt, script)
		script.Pause()
		fmt.Printf("paisong %#v\n", script)
		return nil
	},
	Docstring: "say string	:	outputs string to the dialogue box",
}

//sayf
var DialogueFFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: 1,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptHandler) error {
		//Extract format specifier and values
		fstring := args[0]
		values := make([]interface{}, len(args)-1)
		for i := 0; i < len(values); i++ {
			values[i] = scr.ParseValue(args[i+1])
		}
		txt := fmt.Sprintf(fstring, values...)

		//Update text
		Game.WordHandler.SetText(txt, script)

		script.Pause()
		return nil
	},
	Docstring: "sayf n format ...	:	outputs formatted string to the dialogue box. n is the number of variables + 1. i.e. sayf 2 %s $variable",
}
