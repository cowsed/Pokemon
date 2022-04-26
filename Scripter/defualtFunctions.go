package scripts

import (
	"strconv"
)

//InternalSet acts as the 'cpus' way of setting. It was made its own function because it is used by many other functions
func (scr *ScriptHandler) InternalSet(name, val string) {
	if isVariable(val) {
		val = scr.ParseValue(val)
	}
	scr.memory[name[1:]] = val
}

var SetFunction = ScriptFunction{
	Function: func(args []string, script *Script, scr *ScriptHandler) error {
		name := args[0]
		val := args[1]
		scr.InternalSet(name, val)
		return nil
	},
	NumArguments: 2,
	Docstring: "set  a,  b	:	sets the variable of a to value. b can be a constant value or be prefaced by $ to set a = the value of the variable called b",
}

var AddIFunction = ScriptFunction{
	Function: func(args []string, script *Script, scr *ScriptHandler) error {
		aTok := args[0]
		bTok := args[1]
		c := args[2]

		a := scr.ParseValue(aTok)
		b := scr.ParseValue(bTok)
		aVal, err := strconv.Atoi(a)

		if err != nil {
			return nil
		}
		bVal, err := strconv.Atoi(b)

		if err != nil {
			return nil
		}

		cVal := aVal + bVal

		scr.InternalSet(c, strconv.Itoa(cVal))
		return nil
	},
	NumArguments: 3,
	Docstring: "addI a, b, c	:	adds a+b stores the value to c. a and b can be $ variables or constant values. c must be a variable",
}

var GotoFunc = ScriptFunction{
	NumArguments: 1,
	Function: func(args []string, script *Script, scr *ScriptHandler) error {
		lbl := args[0]
		err := script.Goto(lbl)
		return err
	},
	Docstring: "goto label	:	goes to the label specified",
}

var JmpeFunc = ScriptFunction{
	NumArguments: 3,
	Function: func(args []string, script *Script, scr *ScriptHandler) error {

		lbl := args[0]

		a := scr.ParseValue(args[1])
		b := scr.ParseValue(args[2])

		var err error
		if a == b {
			err = script.Goto(lbl)
		}
		return err
	},
	Docstring: "jmpe lbl, a, b	:	goes to label if a == b. a and b can be variables with $ or constant values",
}

var JmpneFunc = ScriptFunction{
	NumArguments: 3,
	Function: func(args []string, script *Script, scr *ScriptHandler) error {

		lbl := args[0]

		a := scr.ParseValue(args[1])
		b := scr.ParseValue(args[2])

		var err error
		if a != b {
			err = script.Goto(lbl)
		}
		return err
	},
	Docstring: "jmpe lbl, a, b	:	goes to label if a == b. a and b can be variables with $ or constant values",
}
