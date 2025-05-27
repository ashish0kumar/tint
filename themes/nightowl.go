package themes

import "image/color"

var NightOwl = map[string]map[string]color.RGBA{
	"default": {
		"background": hexToRGBA("#011627"), "foreground": hexToRGBA("#d6deeb"),
		"comment": hexToRGBA("#5f7e97"), "blue": hexToRGBA("#82aaff"),
		"green": hexToRGBA("#addb67"), "cyan": hexToRGBA("#7fdbca"),
		"magenta": hexToRGBA("#c792ea"), "red": hexToRGBA("#ef5350"),
		"orange": hexToRGBA("#f78c6c"), "yellow": hexToRGBA("#ffeb95"),
		"caret": hexToRGBA("#80a4c2"), "selection": hexToRGBA("#1d3b53"),
		"lineHighlight": hexToRGBA("#1e293b"), "accent": hexToRGBA("#d6deeb"),
	},
}
