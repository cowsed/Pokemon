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
			"set":  SetFunction,
			"addI": AddIFunction,
			"subI": SubIFunction,
			"mulI": MulIFunction,
			"divI": DivIFunction,

			"addF": AddFFunction,
			"subF": SubFFunction,
			"mulF": MulFFunction,
			"divF": DivFFunction,

			"sqrtF": SqrtFFunction,
			"castI": CastIFunction,

			"goto":  GotoFunc,
			"jmpe":  JmpeFunc,
			"jmpne": JmpneFunc,
			"jmpl":  JmpLessFunc,
			"call":  CallFunc,
			"ret":   RetFunc,
		},
	}
	return &sh
}

func (sh *ScriptEngine) RegisterFunction(name string, sf ScriptFunction) {
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
		action := script.takeToken()

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
			script.End()
			script.Pause()
			stopScript = true

			continue
		}
		if action == "yield" {
			return nil
		}
		//Comment
		if action[0:1] == "#" {
			continue
		}

		f, ok := scr.RegisteredFunctions[action]
		if !ok {
			err := fmt.Errorf("unknown instruction : %s - %v", action, []byte(action))
			return err
		}
		//Take arguments from the code. Handle variable arguments if specified (numArguments < 0)
		argsToGet := f.NumArguments
		if f.NumArguments < 0 {
			numArgs, err := strconv.Atoi(script.takeToken())
			if err != nil {
				return fmt.Errorf("error reading variable arguments on instruction: %s : %v", action, err)
			}
			argsToGet = numArgs
		}
		//Get the arguments
		args := make([]string, argsToGet)
		for i := 0; i < argsToGet; i++ {
			args[i] = script.takeToken()
		}
		//Call the function
		err := f.Function(args, script, scr)
		if err != nil {
			script.Pause()
			return err
		}

		//Check if paused
		if script.paused {
			break
		}

	}
	return nil

}
