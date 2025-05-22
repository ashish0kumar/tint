package themes

import "image/color"

var MonochromePalettes = map[string]map[string]color.RGBA{
	"default": {
		"black":        hexToRGBA("#000000"),
		"darkest-gray": hexToRGBA("#333333"),
		"dark-gray":    hexToRGBA("#666666"),
		"medium-gray":  hexToRGBA("#999999"),
		"light-gray":   hexToRGBA("#CCCCCC"),
		"white":        hexToRGBA("#FFFFFF"),
	},
}
