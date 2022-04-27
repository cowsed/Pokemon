package scripts

import (
	"fmt"
	"strconv"
)

type ScriptEngine struct {
	RegisteredFunctions map[string]ScriptFunction
}

func NewDefaultScriptEngine() *ScriptEngine {
	sh := ScriptEngine{
		RegisteredFunctions: map[string]ScriptFunction{
			"set":   SetFunction,
			"addI":  AddIFunction,
			"goto":  GotoFunc,
			"jmpe":  JmpeFunc,
			"jmpne": JmpneFunc,
			"call":  CallFunc,
			"ret":   RetFunc,
		},
	}
	return &sh
}

func (sh *ScriptEngine) RegisteredFunction(name string, sf ScriptFunction) {
	sh.RegisteredFunctions[name] = sf
}

func (scr *ScriptEngine) ContinueScript(script *Script) error {
	if script.paused {
		return nil
	}

	return scr.ExecuteScript(script)
}

func (scr *ScriptEngine) ExecuteScript(script *Script) error {
	stopScript := false
	for !stopScript {
		//DB print

		//Get the instruction
		action := script.TakeToken()

		//Ignore nonsense
		if len(action) <= 1 {
			continue
		}

		//Dont execute labels
		if isLabel(action) {
			continue
		}

		//End if necessary
		if action == "END" {
			script.Pause()
			stopScript = true
			continue
		}

		f, ok := scr.RegisteredFunctions[action]
		if !ok {
			err := fmt.Errorf("unknown instruction : %s", action)
			return err
		}
		//Take arguments from the code. Handle variable arguments if specified (numArguments < 0)
		argsToGet := f.NumArguments
		if f.NumArguments < 0 {
			numArgs, err := strconv.Atoi(script.TakeToken())
			if err != nil {
				return fmt.Errorf("error reading variable arguments on instruction: %s : %v", action, err)
			}
			argsToGet = numArgs
		}
		//Get the arguments
		args := make([]string, argsToGet)
		for i := 0; i < argsToGet; i++ {
			args[i] = script.TakeToken()
		}
		//Call the function
		f.Function(args, script, scr)

		//Check if paused
		if script.paused {
			break
		}

	}
	return nil

}
