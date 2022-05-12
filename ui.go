package main

import (
	"fmt"
	"io"
	"os"
	ui "pokemon/UI"
	"sort"
	"unicode"

	"github.com/dusk125/pixelui"
	"github.com/faiface/pixel/text"
	"github.com/inkyblackness/imgui-go"
	"golang.org/x/image/font/basicfont"
)

var NoteString string
var ScriptDocsShown = false
var PerformanceShown = true

var showCollisionOverlay = true

func (g *GameStruct) DrawDebugUI() {
	g.ui.NewFrame()
	pc.DrawUI()

	imgui.Begin("Debug")

	imgui.BeginTabBarV("Scripts", 0)

	if imgui.BeginTabItem("Debug Console") {

		s := logger.String()
		imgui.InputTextMultilineV("## Log", &s, imgui.Vec2{X: -1, Y: -1}, imgui.InputTextFlagsReadOnly, nil)

		imgui.EndTabItem()
	}

	if imgui.BeginTabItemV("Render Settings", nil, 0) {
		s := float32(ImageScale)
		imgui.DragFloatV("Scale", &s, 0.05, 0.001, 20, "%.2f", 1)
		imgui.Checkbox("Show Collision Overlay", &showCollisionOverlay)
		ImageScale = float64(s)
		imgui.EndTabItem()
	}

	if imgui.BeginTabItem("Player") {
		g.drawPlayerUi()
		g.player.DrawUI()
		imgui.EndTabItem()
	}

	if imgui.BeginTabItem("Active Entities") {
		drawEntityStatuses(g)
		imgui.EndTabItem()
	}

	if imgui.BeginTabItem("Notepad") {
		imgui.InputTextMultilineV("## Notes", &NoteString, imgui.Vec2{X: -1, Y: -1}, 0, nil)
		imgui.EndTabItem()
	}
	imgui.EndTabBar()

	imgui.End()

	//Script Documentation window
	if ScriptDocsShown {
		imgui.BeginV("Script Documentation", &ScriptDocsShown, 0)
		imgui.Text("Click on a function to view its docstring")
		drawScriptDocs(g)
		imgui.End()
	}
}

var selectedEntity string //String id into map

func drawEntityStatuses(g *GameStruct) {
	//Selected Script Info

	if imgui.Button("Reset all") {
		for _, v := range g.ActiveEntites {
			v.AttachedScript.Restart()
		}
	}

	if selectedEntity != "" {
		entity := g.ActiveEntites[selectedEntity]

		imgui.Text(fmt.Sprintf("%s selected", selectedEntity))
		imgui.Text(entity.AttachedScript.Status())
		

		if imgui.Button("Restart") {
			entity.AttachedScript.Restart()
		}
		imgui.SameLine()
		if imgui.Button("Show Script Docs") {
			ScriptDocsShown = true
		}

		scriptSource := fmt.Sprintf("%v", entity.AttachedScript.MakeHumanReadable(g.ScriptEngine))
		imgui.InputTextMultilineV("## Source", &scriptSource, imgui.Vec2{X: 0, Y: 0}, imgui.InputTextFlagsReadOnly, nil)

		if imgui.TreeNodeV("Internal Rep", imgui.TreeNodeFlagsCollapsingHeader) {
			s := fmt.Sprintf("%#v", entity)
			imgui.PushTextWrapPosV(0)
			imgui.Text(s)
			imgui.PopTextWrapPos()
			imgui.TreePop()
		}

		//Memory
		keys := getKeys(entity.AttachedScript.Memory())

		//imgui.Columns(2, "memory")
		if imgui.TreeNodeV("Memory", imgui.TreeNodeFlagsCollapsingHeader) {

			for _, k := range keys {
				v := entity.AttachedScript.Memory()[k]
				imgui.Separator()
				imgui.Selectable(k + ": ")

				imgui.SameLine()

				imgui.Selectable(v)
			}
			imgui.Separator()
			imgui.TreePop()
		}
	} else {

		imgui.Text("No Entity Selected")
	}

	if imgui.TreeNodeV("Active Entities", imgui.TreeNodeFlagsCollapsingHeader+imgui.TreeNodeFlagsDefaultOpen) {
		//Table of all active scripts
		names := getActiveEntityNames(g)

		imgui.Separator()
		for _, name := range names {
			if imgui.Selectable(name) {
				selectedEntity = name
			}
			imgui.Separator()
		}
		imgui.TreePop()
	}
}

func (g *GameStruct) drawPlayerUi() {
	imgui.Text(fmt.Sprintf("%.2f x %.2f", g.player.x, g.player.y))
}

func getKeys(mem map[string]string) []string {
	ks := make([]string, len(mem))
	i := 0
	for k := range mem {
		ks[i] = k
		i++
	}
	sort.Strings(ks)
	return ks
}

func getActiveEntityNames(g *GameStruct) []string {
	names := make([]string, len(g.ActiveEntites))
	i := 0
	for k := range g.ActiveEntites {
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
	Game.WordHandler = &ui.DialogueHandler{
		WaitingForConfirmation: false,
		ListedText:             "",
		Active:                 false,
		Drawer:                 ui.NewTextDrawer(g.atlas, g.win.Bounds().Size().X, g.win.Bounds().Size().Y),
		UIHandler:              g.ui,
		OnClose:                func() {},
	}
}

func (g *GameStruct) InitializeUI() {
	LoadNote()
	g.ui = pixelui.NewUI(g.win, pixelui.NO_DEFAULT_FONT)

	g.ui.AddTTFFont("Resources/FreeMono.ttf", 24)

	imgui.CurrentStyle().ScaleAllSizes(2)

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
