package main

import (
	"fmt"
	scripts "pokemon/Scripter"
)

var DebugLogFunction scripts.ScriptFunction = scripts.ScriptFunction{
	NumArguments: 1,
	Function: func(args []string, script *scripts.Script, scr *scripts.ScriptHandler) error {
		Game.logger.Write([]byte(args[0] + "\n"))
		return nil
	},
	Docstring: "Logs the specified text to the game console",
}

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
