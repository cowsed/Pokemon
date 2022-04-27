package scripts

import (
	"fmt"
	"strconv"
	"strings"
)

//Syntax
/*
	label:			   	-> a string followed by: will be treated as a label
	set  a,  b 		   	-> sets the variable of a to value. b can be a constant value or be prefaced by $ to set a = the value of the variable called b
	say  a    		   	-> prints out the string a
	sayf n, f, ...  	-> does printf(f, ...) where ... is any number of variables or values and n is the length of ...
	ask  n, q, ...  	-> asks the user for input. prints q, then asks to pick from one of ... where n is the length of ... (oregon trail style entering for now. press one for option 1)
	goto lbl			-> goes to label
	jmpe lbl, a, b		-> goes to label if a == b. a and b can be variables with $ or constant values
	jmpne lbl, a, b		-> goes to label if a != b. a and b can be variables with $ or constant values
	addI a, b, c		-> adds a+b stores the value to c. a and b can be $ variables or constant values. c must be a variable
*/

type ScriptFunction struct {
	NumArguments int //The number of tokens to read and give to the function. If negative a variable argument is assumed. The next token should be an int saying how many it needs

	Function  func(args []string, script *Script, scr *ScriptHandler) error //What to call when this function is seen
	Docstring string
}

func (sh *ScriptHandler) RegisteredFunction(name string, sf ScriptFunction) {
	sh.RegisteredFunctions[name] = sf
}

type ScriptHandler struct {
	RegisteredFunctions map[string]ScriptFunction
	PrintFunc           func(s string)
	memory              map[string]string
}
type Script struct {
	src    []string
	index  int
	labels map[string]int
	paused bool
}

var labelChar = ":"
var variableChar = "$"

func NewDefualtScriptEngine() *ScriptHandler {
	sh := ScriptHandler{
		PrintFunc: nil,
		memory:    map[string]string{},
		RegisteredFunctions: map[string]ScriptFunction{
			"set":   SetFunction,
			"addI":  AddIFunction,
			"goto":  GotoFunc,
			"jmpe":  JmpeFunc,
			"jmpne": JmpneFunc,
		},
	}
	return &sh
}

func isLabel(w string) bool {
	res := w[len(w)-1:] == labelChar
	return res
}
func isVariable(w string) bool {
	res := w[0:1] == variableChar
	return res

}
func NewScript(src string) *Script {

	srcTokens := Split(src)
	s := Script{
		src:    srcTokens,
		index:  0,
		labels: map[string]int{},
		paused: false,
	}
	for i, w := range srcTokens {
		if len(w) < 1 {
			continue
		}
		if isLabel(w) {
			l := w[:len(w)-1]
			s.labels[l] = i
		}
	}

	return &s
}
func (s *Script) Resume() {
	s.paused = false
}
func (s *Script) Pause() {
	s.paused = true
}
func (s *Script) Restart() {
	s.paused = false
	s.index = 0
}
func (s *Script) Status() string {
	str := "running"
	if s.paused {
		str = "Waiting For Events"
	}
	return fmt.Sprintf(str+" - PC: %v", s.index)
}

func (s *Script) MakeHumanReadable(sh *ScriptHandler) string {
	sourceText := ""
	//Read through tokens, get arguments specified by the function. then make new line
	index := 0
	for index < len(s.src) {
		line := ""
		functionName := s.src[index]
		index++

		if functionName == "" {
			continue
		}

		line += functionName
		if isLabel(functionName) {
			sourceText += line + "\n"
			continue
		}
		function, ok := sh.RegisteredFunctions[functionName]
		if !ok {
			sourceText = "Unrecognized function: '" + functionName + "'"
			break
		}
		numArgs := function.NumArguments

		if function.NumArguments < 0 {
			varArgsNumString := s.src[index]
			index++
			numArgs, _ = strconv.Atoi(varArgsNumString)
		}

		//Take the arguments off the top
		for i := 0; i < numArgs; i++ {
			arg := s.src[index]
			//if the token was originally enclosed by quotes, it would have spaces
			if strings.Contains(arg, " ") {
				arg = "\"" + arg + "\""
			}
			line += " " + arg
			index++
		}
		sourceText += line + "\n"
	}
	return sourceText
}

func (scr *ScriptHandler) ContinueScript(script *Script) error {
	if script.paused {
		return nil
	}

	return scr.ExecuteScript(script)
}

func (scr *ScriptHandler) ExecuteScript(script *Script) error {
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
func (s *Script) Goto(lbl string) error {
	val, ok := s.labels[lbl]
	if !ok {
		return fmt.Errorf("no label called %s", lbl)
	}
	s.index = val
	return nil
}

func (s *Script) TakeToken() string {
	if s.index >= len(s.src) {
		return "END"
	}
	str := s.src[s.index]
	s.index++
	return str
}

func (scr *ScriptHandler) ParseValue(name string) string {
	if isVariable(name) { //Remove% if still there
		return scr.memory[name[1:]]
	}
	return name
}

//Splits the input string into tokens
func Split(instring string) []string {
	inside_string := false
	tokens := make([]string, 0, 30)
	currentToken := ""
	ignoreNext := false
	for index, char := range instring {
		if ignoreNext {
			ignoreNext = false
			continue
		}
		schar := string(char)
		if inside_string {

			//Another " end that string
			if schar == "\"" {
				inside_string = false
			} else if schar == "\\" {
				//Handle specials
				nextChar := string(instring[index+1])
				ignoreNext = true
				switch nextChar {
				case "n":
					currentToken += "\n"
				case "t":
					currentToken += "\t"

				}

			} else {
				currentToken += schar
			}
			continue
		}

		if schar == " " || schar == "\n" {
			//Ignore double white spaces. dont make an empty token
			if currentToken != "" {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}
		} else if schar == "\"" {
			inside_string = !inside_string
		} else {
			currentToken += schar
		}
	}
	tokens = append(tokens, currentToken)
	return tokens
}
