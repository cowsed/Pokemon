package main

import (
	"io"
	"os"
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

var SelectedFunction = ""

func drawScriptDocs(g *GameStruct) {
	names := g.ScriptEngine.FunctionNames()
	for _, n := range names {
		var Selected = SelectedFunction == n
		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{.75, .75, .75, 1})

		if Selected {
			imgui.PopStyleColor()

			imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{1, 1, 1, 1})

		}
		if imgui.Selectable(n) {
			SelectedFunction = n
		}
		imgui.PopStyleColor()
	}
	descText := g.ScriptEngine.GetDocstring(SelectedFunction)
	if descText == "" {
		descText = "No Documentation Available"
	}
	//imgui.InputTextMultilineV("## Documentation", &descText, imgui.Vec2{X: -1, Y: -1}, 0, imgui.)
	imgui.PushTextWrapPosV(0)
	imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
	imgui.Text(descText)
	imgui.PopStyleColor()
	imgui.PopTextWrapPos()
}

func (g *GameStruct) InitializeGameUI() {
	g.atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII, text.RangeTable(unicode.Latin))

	//Setup Bottom of the window dialogue text thing
	t := NewTextDrawer(g.atlas, g.win.Bounds().Size().X, g.win.Bounds().Size().Y)
	Game.DialogueHandler = t
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
