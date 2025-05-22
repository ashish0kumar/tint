package themes

import "image/color"

var NordPalettes = map[string]map[string]color.RGBA{
	"default": {
		// Polar Night
		"nord0": hexToRGBA("#2E3440"),
		"nord1": hexToRGBA("#3B4252"),
		"nord2": hexToRGBA("#434C5E"),
		"nord3": hexToRGBA("#4C566A"),

		// Snow Storm
		"nord4": hexToRGBA("#D8DEE9"),
		"nord5": hexToRGBA("#E5E9F0"),
		"nord6": hexToRGBA("#ECEFF4"),

		// Frost
		"nord7":  hexToRGBA("#8FBCBB"),
		"nord8":  hexToRGBA("#88C0D0"),
		"nord9":  hexToRGBA("#81A1C1"),
		"nord10": hexToRGBA("#5E81AC"),

		// Aurora
		"nord11": hexToRGBA("#BF616A"),
		"nord12": hexToRGBA("#D08770"),
		"nord13": hexToRGBA("#EBCB8B"),
		"nord14": hexToRGBA("#A3BE8C"),
		"nord15": hexToRGBA("#B48EAD"),
	},
}
