package themes

import "image/color"

var EverforestPalettes = map[string]map[string]color.RGBA{
	"dark": {
		"background":  hexToRGBA("#2b3339"),
		"foreground":  hexToRGBA("#d3c6aa"),
		"currentLine": hexToRGBA("#343f44"),
		"comment":     hexToRGBA("#7a8478"),
		"red":         hexToRGBA("#e67e80"),
		"orange":      hexToRGBA("#e69875"),
		"yellow":      hexToRGBA("#dbbc7f"),
		"green":       hexToRGBA("#a7c080"),
		"teal":        hexToRGBA("#83c092"),
		"blue":        hexToRGBA("#7fbbb3"),
		"purple":      hexToRGBA("#d699b6"),
		"magenta":     hexToRGBA("#d699b6"),
		"cyan":        hexToRGBA("#83c092"),
		"white":       hexToRGBA("#d3c6aa"),
		"black":       hexToRGBA("#2b3339"),
		"gray":        hexToRGBA("#7a8478"),
	},
	"light": {
		"background":  hexToRGBA("#fdf6e3"),
		"foreground":  hexToRGBA("#534d40"),
		"currentLine": hexToRGBA("#f5eddc"),
		"comment":     hexToRGBA("#a6b0a0"),
		"red":         hexToRGBA("#e66e77"),
		"orange":      hexToRGBA("#e6914f"),
		"yellow":      hexToRGBA("#d5a84b"),
		"green":       hexToRGBA("#8da13b"),
		"teal":        hexToRGBA("#3a948c"),
		"blue":        hexToRGBA("#3a948c"),
		"purple":      hexToRGBA("#df627d"),
		"magenta":     hexToRGBA("#df627d"),
		"cyan":        hexToRGBA("#3a948c"),
		"white":       hexToRGBA("#534d40"),
		"black":       hexToRGBA("#fdf6e3"),
		"gray":        hexToRGBA("#a6b0a0"),
	},
}
