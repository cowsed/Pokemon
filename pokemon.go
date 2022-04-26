package main

type PokemonType string

type Species struct {
	Number int

	Name string

	DexDescription string

	Type1 PokemonType
	Type2 PokemonType

	Height float64
	Wieght float64

	LevelToEvolve uint8

	GenderRatio float64 //percentage males. .5 -> half and half. .75 -> .75 male, .25 female

	Ability       string
	HiddenAbility string

	HatchtimeLow, HatchtimeHigh int

	BaseExperienceYield int

	EggGroups []string

	LearnableTMS   []int
	LearnableMoves map[int][]string

	CaptureRate int
	IsLegendary bool
}

var TestFeraligatr = Species{
	Number:              160,
	Name:                "Feraligatr",
	DexDescription:      "When it bites with its massive and powerful jaws, it shakes its head and savagely tears its victim up.",
	Type1:               "water",
	Type2:               "none",
	Height:              2.3,
	Wieght:              45,
	LevelToEvolve:       255,
	GenderRatio:         .875,
	Ability:             "torrent",
	HiddenAbility:       "sheer force",
	HatchtimeLow:        5140,
	HatchtimeHigh:       5396,
	BaseExperienceYield: 239,
	EggGroups:           []string{"monster", "water1"},
	LearnableTMS:        []int{1, 2, 3, 5, 7, 10, 13, 14, 15, 17, 18, 23, 26, 28, 31, 32, 39, 40, 42, 44, 45, 49, 52, 56, 58, 59, 65, 68, 72, 75, 80, 82, 83, 87, 90, 93, 95, 96, 98, 100},
	LearnableMoves: map[int][]string{
		1:  {"Agility", "Scratch", "Leer", "Water Gun", "Mud-Slap"},
		6:  {"Water Gun"},
		8:  {"Mud-Slap"},
		13: {"Bite"},
		15: {"Scary Face"},
		21: {"Ice Fang"},
		24: {"Flail"},
		32: {"Crunch"},
		37: {"Low Kick"},
		45: {"Slash"},
		50: {"Screech"},
		58: {"Thrash"},
		63: {"Aqua Tail"},
		71: {"Superpower"},
		76: {"Hydro Punp"},
	},
}
