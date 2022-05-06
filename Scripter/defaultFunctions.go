package scripts

import (
	"math"
	"strconv"
)

//internalSet acts as the 'cpus' way of setting. It was made its own function because it is used by many other functions
func (s *Script) internalSet(name, val string) {
	if isVariable(val) {
		val = s.ParseValue(val)
	}
	s.memory[name[1:]] = val
}

var SetFunction = ScriptFunction{
	Function: func(args []string, script *Script, scr *ScriptEngine) error {
		name := args[0]
		val := args[1]
		script.internalSet(name, val)
		return nil
	},
	NumArguments: 2,
	Docstring: "set  a,  b	:	sets the variable of a to value. b can be a constant value or be prefaced by $ to set a = the value of the variable called b",
}

var AddIFunction = ScriptFunction{
	Function: func(args []string, script *Script, scr *ScriptEngine) error {
		aTok := args[0]
		bTok := args[1]
		c := args[2]

		a := script.ParseValue(aTok)
		b := script.ParseValue(bTok)
		aVal, err := strconv.Atoi(a)

		if err != nil {
			return err
		}
		bVal, err := strconv.Atoi(b)

		if err != nil {
			return err
		}

		cVal := aVal + bVal

		script.internalSet(c, strconv.Itoa(cVal))
		return nil
	},
	NumArguments: 3,
	Docstring: "addI a, b, c	:	adds a+b stores the value to c. a and b can be $ variables or constant values. c must be a variable",
}

var SubIFunction = ScriptFunction{
	Function: func(args []string, script *Script, scr *ScriptEngine) error {
		aTok := args[0]
		bTok := args[1]
		c := args[2]

		a := script.ParseValue(aTok)
		b := script.ParseValue(bTok)
		aVal, err := strconv.Atoi(a)

		if err != nil {
			return err
		}
		bVal, err := strconv.Atoi(b)

		if err != nil {
			return err
		}

		cVal := aVal - bVal

		script.internalSet(c, strconv.Itoa(cVal))
		return nil
	},
	NumArguments: 3,
	Docstring: "subI a, b, c	:	does a-b stores the value to c. a and b can be $ variables or constant values. c must be a variable",
}
var MulIFunction = ScriptFunction{
	Function: func(args []string, script *Script, scr *ScriptEngine) error {
		aTok := args[0]
		bTok := args[1]
		c := args[2]

		a := script.ParseValue(aTok)
		b := script.ParseValue(bTok)
		aVal, err := strconv.Atoi(a)

		if err != nil {
			return err
		}
		bVal, err := strconv.Atoi(b)

		if err != nil {
			return err
		}

		cVal := aVal * bVal

		script.internalSet(c, strconv.Itoa(cVal))
		return nil
	},
	NumArguments: 3,
	Docstring: "mulI a, b, c	:	multiplies a*b stores the value to c. a and b can be $ variables or constant values. c must be a variable",
}
var DivIFunction = ScriptFunction{
	Function: func(args []string, script *Script, scr *ScriptEngine) error {
		aTok := args[0]
		bTok := args[1]
		c := args[2]

		a := script.ParseValue(aTok)
		b := script.ParseValue(bTok)
		aVal, err := strconv.Atoi(a)

		if err != nil {
			return err
		}
		bVal, err := strconv.Atoi(b)

		if err != nil {
			return err
		}

		cVal := aVal / bVal

		script.internalSet(c, strconv.Itoa(cVal))
		return nil
	},
	NumArguments: 3,
	Docstring: "divI a, b, c	:	divides a/b stores the value to c. a and b can be $ variables or constant values. c must be a variable",
}

var AddFFunction = ScriptFunction{
	Function: func(args []string, script *Script, scr *ScriptEngine) error {
		aTok := args[0]
		bTok := args[1]
		c := args[2]

		a := script.ParseValue(aTok)
		b := script.ParseValue(bTok)
		aVal, err := strconv.ParseFloat(a, 64)

		if err != nil {
			return err
		}
		bVal, err := strconv.ParseFloat(b, 64)

		if err != nil {
			return err
		}

		cVal := aVal + bVal

		script.internalSet(c, strconv.FormatFloat(cVal, 'g', 8, 64))
		return nil
	},
	NumArguments: 3,
	Docstring: "addF a, b, c	:	adds a+b stores the value to c. a and b can be $ variables or constant values. c must be a variable",
}
var SubFFunction = ScriptFunction{
	Function: func(args []string, script *Script, scr *ScriptEngine) error {
		aTok := args[0]
		bTok := args[1]
		c := args[2]

		a := script.ParseValue(aTok)
		b := script.ParseValue(bTok)
		aVal, err := strconv.ParseFloat(a, 64)

		if err != nil {
			return err
		}
		bVal, err := strconv.ParseFloat(b, 64)

		if err != nil {
			return err
		}

		cVal := aVal - bVal

		script.internalSet(c, strconv.FormatFloat(cVal, 'g', 8, 64))
		return nil
	},
	NumArguments: 3,
	Docstring: "subF a, b, c	:	adds a+b stores the value to c. a and b can be $ variables or constant values. c must be a variable",
}

var MulFFunction = ScriptFunction{
	Function: func(args []string, script *Script, scr *ScriptEngine) error {
		aTok := args[0]
		bTok := args[1]
		c := args[2]

		a := script.ParseValue(aTok)
		b := script.ParseValue(bTok)
		aVal, err := strconv.ParseFloat(a, 64)

		if err != nil {
			return err
		}
		bVal, err := strconv.ParseFloat(b, 64)

		if err != nil {
			return err
		}

		cVal := aVal * bVal

		script.internalSet(c, strconv.FormatFloat(cVal, 'g', 8, 64))
		return nil
	},
	NumArguments: 3,
	Docstring: "mulF a, b, c	:	adds a+b stores the value to c. a and b can be $ variables or constant values. c must be a variable",
}

var DivFFunction = ScriptFunction{
	Function: func(args []string, script *Script, scr *ScriptEngine) error {
		aTok := args[0]
		bTok := args[1]
		c := args[2]

		a := script.ParseValue(aTok)
		b := script.ParseValue(bTok)
		aVal, err := strconv.ParseFloat(a, 64)

		if err != nil {
			return err
		}
		bVal, err := strconv.ParseFloat(b, 64)

		if err != nil {
			return err
		}

		cVal := aVal / bVal

		script.internalSet(c, strconv.FormatFloat(cVal, 'g', 8, 64))
		return nil
	},
	NumArguments: 3,
	Docstring: "divF a, b, c	:	adds a+b stores the value to c. a and b can be $ variables or constant values. c must be a variable",
}

var SqrtFFunction = ScriptFunction{
	Function: func(args []string, script *Script, scr *ScriptEngine) error {
		aTok := args[0]
		c := args[1]

		a := script.ParseValue(aTok)
		aVal, err := strconv.ParseFloat(a, 64)

		if err != nil {
			return err
		}

		cVal := math.Sqrt(aVal)

		script.internalSet(c, strconv.FormatFloat(cVal, 'g', 8, 64))
		return nil
	},
	NumArguments: 2,
	Docstring: "sqrtF a, b	:	sqrts(a) stores the value to b. a can be $ variables or constant values. b must be a variable",
}

var CastIFunction = ScriptFunction{
	Function: func(args []string, script *Script, scr *ScriptEngine) error {
		aTok := args[0]
		a := script.ParseValue(aTok)
		f, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return err
		}
		i := int(f)
		bTok := args[1] //location to store location
		script.internalSet(bTok, strconv.FormatInt(int64(i), 10))

		return nil
	},
	NumArguments: 2,
	Docstring: "castI a,b	:	Casts the value a to an int and stores it to b. a can be a variable or literal. b must be a variable name",
}

var GotoFunc = ScriptFunction{
	NumArguments: 1,
	Function: func(args []string, script *Script, scr *ScriptEngine) error {
		lbl := args[0]
		err := script.Goto(lbl)
		return err
	},
	Docstring: "goto label	:	goes to the label specified",
}

var JmpeFunc = ScriptFunction{
	NumArguments: 3,
	Function: func(args []string, script *Script, scr *ScriptEngine) error {

		lbl := args[0]

		a := script.ParseValue(args[1])
		b := script.ParseValue(args[2])

		var err error
		if a == b {
			err = script.Goto(lbl)
		}
		return err
	},
	Docstring: "jmpe lbl, a, b	:	goes to label if a == b. a and b can be variables with $ or constant values",
}

var JmpLessFunc = ScriptFunction{
	NumArguments: 3,
	Function: func(args []string, script *Script, scr *ScriptEngine) error {

		lbl := args[0]

		a := script.ParseValue(args[1])
		b := script.ParseValue(args[2])

		aVal, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return err
		}
		bVal, err := strconv.ParseFloat(b, 64)
		if err != nil {
			return err
		}

		if aVal < bVal {
			err = script.Goto(lbl)
		}
		return err
	},
	Docstring: "jmpl lbl, a, b	:	goes to label if a < b. a and b can be variables with $ or constant values of any numerical value",
}

var JmpneFunc = ScriptFunction{
	NumArguments: 3,
	Function: func(args []string, script *Script, scr *ScriptEngine) error {

		lbl := args[0]

		a := script.ParseValue(args[1])
		b := script.ParseValue(args[2])

		var err error
		if a != b {
			err = script.Goto(lbl)
		}
		return err
	},
	Docstring: "jmpe lbl, a, b	:	goes to label if a == b. a and b can be variables with $ or constant values",
}

var CallFunc = ScriptFunction{
	NumArguments: 1,
	Function: func(args []string, script *Script, scr *ScriptEngine) error {
		lbl := args[0]
		script.stack.Push(script.index)
		script.Goto(lbl)
		return nil
	},
	Docstring: "call label	:	pushes the current program counter goes to the label specified. when a corresponding ret function is called, execution will resume here. Call subroutine",
}

var RetFunc = ScriptFunction{
	NumArguments: 0,
	Function: func(args []string, script *Script, scr *ScriptEngine) error {
		pos := script.stack.Pop()
		script.index = pos
		return nil
	},
	Docstring: "ret 	:	returns to the position in the program given by call last time. Return from subroutine",
}
