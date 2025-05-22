package themes

import "image/color"

var DraculaPalettes = map[string]map[string]color.RGBA{
	"default": {
		"background":   hexToRGBA("#282a36"),
		"current-line": hexToRGBA("#44475a"),
		"foreground":   hexToRGBA("#f8f8f2"),
		"comment":      hexToRGBA("#6272a4"),
		"cyan":         hexToRGBA("#8be9fd"),
		"green":        hexToRGBA("#50fa7b"),
		"orange":       hexToRGBA("#ffb86c"),
		"pink":         hexToRGBA("#ff79c6"),
		"purple":       hexToRGBA("#bd93f9"),
		"red":          hexToRGBA("#ff5555"),
		"yellow":       hexToRGBA("#f1fa8c"),
	},
}
