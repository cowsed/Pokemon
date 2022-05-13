package scripts

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var labelChar = ":"
var variableChar = "$"

//A consistent way for a script to interact with an outside system. added into the Script Engine by RegisterFunction
type ScriptFunction struct {
	NumArguments int //The number of tokens to read and give to the function. If negative a variable argument is assumed. The next token should be an int saying how many it needs

	Function  func(args []string, script *Script, scr *ScriptEngine) error //What to call when this function is seen
	Docstring string
}

type Script struct {
	src    []string
	index  int
	labels map[string]int
	paused bool
	ended  bool
	memory map[string]string
	stack  *Stack
}

func NewScriptFromFile(filepath string) *Script {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error loading script from file:", err)
	}
	bs, _ := io.ReadAll(f)
	f.Close()

	src := string(bs)
	return NewScript(src)
}

func NewScript(src string) *Script {
	src = string(NormalizeNewlines([]byte(src)))
	srcTokens := Split(src)
	s := Script{
		src:    srcTokens,
		index:  0,
		labels: map[string]int{},
		paused: false,
		ended:  false,
		memory: map[string]string{},
		stack:  NewStack(),
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
func (s *Script) ClearMemory() {
	s.memory = map[string]string{}
}

func (s *Script) Memory() map[string]string {
	return s.memory
}

func (s *Script) Resume() {
	s.paused = false
}
func (s *Script) Pause() {
	s.paused = true
}

func (s *Script) Restart() {
	s.paused = false
	s.ended = false
	s.index = 0
}
func (s *Script) End() {
	s.ended = true
}

func (s *Script) Backup(numTokens int) {
	s.index -= numTokens
	if s.index < 0 {
		s.index = 0
	}
}

func (s *Script) SetMemory(key, value string) {
	s.memory[key] = value
}
func (s *Script) GetMemory(key string) string {
	return s.memory[key]
}

func (s *Script) Status() string {
	str := "running"
	if s.paused {
		str = "Waiting For Events"
	}
	if s.ended {
		str += "  -  finished"
	}
	return fmt.Sprintf(str+" - PC: %v/%v", s.index, len(s.src))
}
func (s *Script) HasLabel(lbl string) bool {
	_, inmap := s.labels[lbl]
	return inmap
}

func (s *Script) MakeHumanReadable(sh *ScriptEngine) string {
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
		if functionName == "END" || functionName == "yield" || functionName[0:1] == "#" {
			sourceText += functionName + "\n"
			continue
		}
		line += functionName
		if isLabel(functionName) {
			sourceText += line + "\n"
			continue
		}
		function, ok := sh.RegisteredFunctions[functionName]
		if !ok {
			sourceText = "Unrecognized function: '" + functionName + "' " + fmt.Sprint([]byte(functionName))
			fmt.Println(s.src)
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

func (s *Script) Goto(lbl string) error {
	val, ok := s.labels[lbl]
	if !ok {
		return fmt.Errorf("no label called %s", lbl)
	}
	s.index = val
	if s.ended {
		s.ended = false
		s.paused = false
	}

	return nil
}

func (s *Script) Call(lbl string) error {
	val, ok := s.labels[lbl]
	if !ok {
		return fmt.Errorf("no label called %s", lbl)
	}
	s.stack.Push(s.index)
	s.index = val

	if s.ended {
		s.ended = false
		s.paused = false
	}

	return nil
}

func (s *Script) takeToken() string {
	if s.index >= len(s.src) {
		return "END"
	}
	str := s.src[s.index]
	s.index++
	return str
}

func (s *Script) ParseValue(name string) string {
	if isVariable(name) { //Remove$ if still there
		return s.memory[name[1:]]
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

func isVariable(w string) bool {
	res := w[0:1] == variableChar
	return res

}

func isLabel(w string) bool {
	res := w[len(w)-1:] == labelChar
	return res
}

func NormalizeNewlines(d []byte) []byte {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)
	return d
}
