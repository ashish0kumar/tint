package themes

import "image/color"

var Ayu = map[string]map[string]color.RGBA{
	"default": {
		"background": hexToRGBA("#0f1419"), "foreground": hexToRGBA("#b3b1ad"),
		"comment": hexToRGBA("#5c6773"), "cyan": hexToRGBA("#36a3d9"),
		"blue": hexToRGBA("#39bae6"), "purple": hexToRGBA("#c296eb"),
		"green": hexToRGBA("#aad94c"), "red": hexToRGBA("#f07178"),
		"orange": hexToRGBA("#f29718"), "yellow": hexToRGBA("#ffb454"),
		"selection": hexToRGBA("#253340"), "lineHighlight": hexToRGBA("#14191f"),
		"caret": hexToRGBA("#ffcc66"), "accent": hexToRGBA("#5ccfe6"),
	},

	"light": {
		"background": hexToRGBA("#fafafa"), "foreground": hexToRGBA("#5c6773"),
		"comment": hexToRGBA("#abb0b6"), "cyan": hexToRGBA("#55b4d4"),
		"blue": hexToRGBA("#36a3d9"), "purple": hexToRGBA("#a37acc"),
		"green": hexToRGBA("#86b300"), "red": hexToRGBA("#f07171"),
		"orange": hexToRGBA("#f29718"), "yellow": hexToRGBA("#ff9940"),
		"selection": hexToRGBA("#e5e5e6"), "lineHighlight": hexToRGBA("#f3f4f5"),
		"caret": hexToRGBA("#ff9940"), "accent": hexToRGBA("#36a3d9"),
	},

	"mirage": {
		"background": hexToRGBA("#1f2430"), "foreground": hexToRGBA("#cbccc6"),
		"comment": hexToRGBA("#5c6773"), "cyan": hexToRGBA("#95e6cb"),
		"blue": hexToRGBA("#91ddff"), "purple": hexToRGBA("#d4bfff"),
		"green": hexToRGBA("#bae67e"), "red": hexToRGBA("#f28779"),
		"orange": hexToRGBA("#ffcc66"), "yellow": hexToRGBA("#ffd580"),
		"selection": hexToRGBA("#33415e"), "lineHighlight": hexToRGBA("#2d3347"),
		"caret": hexToRGBA("#ffcc66"), "accent": hexToRGBA("#5ccfe6"),
	},

	"dark": {
		"background": hexToRGBA("#0f1419"), "foreground": hexToRGBA("#b3b1ad"),
		"comment": hexToRGBA("#5c6773"), "cyan": hexToRGBA("#36a3d9"),
		"blue": hexToRGBA("#39bae6"), "purple": hexToRGBA("#c296eb"),
		"green": hexToRGBA("#aad94c"), "red": hexToRGBA("#f07178"),
		"orange": hexToRGBA("#f29718"), "yellow": hexToRGBA("#ffb454"),
		"selection": hexToRGBA("#253340"), "lineHighlight": hexToRGBA("#14191f"),
		"caret": hexToRGBA("#ffcc66"), "accent": hexToRGBA("#5ccfe6"),
	},
}
