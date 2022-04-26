package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"unicode"

	"github.com/dusk125/pixelui"
	"github.com/faiface/pixel/text"
	"github.com/inkyblackness/imgui-go"
	"golang.org/x/image/font/basicfont"
)

var NoteString string

func (g *GameStruct) DrawDebugUI() {
	g.ui.NewFrame()

	imgui.Begin("Debug")

	imgui.BeginTabBar("Scripts")

	if imgui.BeginTabItem("Debug Console") {

		s := g.logger.String()
		imgui.InputTextMultilineV("## Log", &s, imgui.Vec2{X: -1, Y: -1}, imgui.InputTextFlagsReadOnly, nil)

		imgui.EndTabItem()
	}
	if imgui.BeginTabItem("Active Scripts") {
		drawScriptStatuses(g)
		imgui.EndTabItem()
	}
	if imgui.BeginTabItem("Script Documentation") {
		drawScriptDocs(g)
		imgui.EndTabItem()
	}

	if imgui.BeginTabItem("Notepad") {
		imgui.InputTextMultilineV("## Notes", &NoteString, imgui.Vec2{X: -1, Y: -1}, 0, nil)
		imgui.EndTabItem()
	}
	imgui.EndTabBar()

	imgui.End()

}

var selectedScript string //String id into map
func drawScriptStatuses(g *GameStruct) {
	//Selected Script Info
	if selectedScript != "" {
		script := g.ActiveScripts[selectedScript]

		//Overview of script
		imgui.Text(fmt.Sprintf("%s selected", selectedScript))
		imgui.Text(script.Status())

		if imgui.Button("Restart") {
			script.Restart()
		}

		//Show source code of script
		scriptSource := fmt.Sprintf("%v", script.MakeHumanReadable(g.ScriptEngine))
		imgui.InputTextMultilineV("## Source", &scriptSource, imgui.Vec2{0, 0}, imgui.InputTextFlagsReadOnly, nil)

	}

	//Table of all active scripts
	names := getActiveScriptNames(g)

	imgui.Text("All Active Scripts:")
	imgui.Separator()
	for _, name := range names {
		if imgui.Selectable(name) {
			selectedScript = name
		}
		imgui.Separator()

	}

}
func getActiveScriptNames(g *GameStruct) []string {
	names := make([]string, len(g.ActiveScripts))
	i := 0
	for k := range g.ActiveScripts {
		names[i] = k
		i++
	}
	sort.Strings(names)
	return names
}

var SelectedFunctionForDocs = ""

func drawScriptDocs(g *GameStruct) {
	names := g.ScriptEngine.FunctionNames()
	for _, n := range names {
		var Selected = SelectedFunctionForDocs == n
		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: .75, Y: .75, Z: .75, W: 1})

		if Selected {
			imgui.PopStyleColor()

			imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 1, Z: 1, W: 1})

		}
		if imgui.Selectable(n) {
			SelectedFunctionForDocs = n
		}
		imgui.PopStyleColor()
	}

	//Get the documentation text or supply defualt
	descText := g.ScriptEngine.GetDocstring(SelectedFunctionForDocs)
	if descText == "" {
		descText = "No Documentation Available"
	}

	//Until an approach to line wrapping, this is it
	imgui.PushTextWrapPosV(0)
	imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
	imgui.Text(descText)
	imgui.PopStyleColor()
	imgui.PopTextWrapPos()
}

func (g *GameStruct) InitializeGameUI() {
	g.atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII, text.RangeTable(unicode.Latin))

	//Setup Bottom of the window dialogue text thing
	Game.WordHandler = &DialogueHandler{
		WaitingForConfirmation: true,
		ListedText:             "",
		Active:                 true,
		drawer:                 NewTextDrawer(g.atlas, g.win.Bounds().Size().X, g.win.Bounds().Size().Y),
	}
}

func (g *GameStruct) InitializeUI() {
	LoadNote()
	g.ui = pixelui.NewUI(g.win, pixelui.NO_DEFAULT_FONT)

	g.ui.AddTTFFont("Resources/FreeMono.ttf", 24)

	imgui.CurrentStyle().ScaleAllSizes(2)

	Game.logger = &Logger{}

}

//Loads and unloads the note in the notepad
func LoadNote() {
	f, err := os.Open("Resources/note.txt")
	check(err)
	bs, _ := io.ReadAll(f)
	NoteString = string(bs)
}
func SaveNote() {
	f, err := os.Create("Resources/note.txt")
	check(err)
	f.Write([]byte(NoteString))
	f.Close()

}
