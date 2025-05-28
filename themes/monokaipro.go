package themes

import "image/color"

var MonokaiPro = map[string]map[string]color.RGBA{
	"default": {
		"background": hexToRGBA("#2d2a2e"), "foreground": hexToRGBA("#fcfcfa"),
		"comment": hexToRGBA("#727072"), "red": hexToRGBA("#ff6188"),
		"orange": hexToRGBA("#fc9867"), "yellow": hexToRGBA("#ffd866"),
		"green": hexToRGBA("#a9dc76"), "cyan": hexToRGBA("#78dce8"),
		"blue": hexToRGBA("#78dce8"), "purple": hexToRGBA("#ab9df2"),
	},

	"classic": {
		"background": hexToRGBA("#2d2a2e"), "foreground": hexToRGBA("#fcfcfa"),
		"comment": hexToRGBA("#727072"), "red": hexToRGBA("#ff6188"),
		"orange": hexToRGBA("#fc9867"), "yellow": hexToRGBA("#ffd866"),
		"green": hexToRGBA("#a9dc76"), "cyan": hexToRGBA("#78dce8"),
		"blue": hexToRGBA("#78dce8"), "purple": hexToRGBA("#ab9df2"),
	},

	"spectrum": {
		"background": hexToRGBA("#2d2a2e"), "foreground": hexToRGBA("#fcfcfa"),
		"comment": hexToRGBA("#727072"), "red": hexToRGBA("#ff6188"),
		"orange": hexToRGBA("#fc9867"), "yellow": hexToRGBA("#ffd866"),
		"green": hexToRGBA("#a9dc76"), "cyan": hexToRGBA("#78dce8"),
		"blue": hexToRGBA("#78dce8"), "purple": hexToRGBA("#ab9df2"),
	},

	"octagon": {
		"background": hexToRGBA("#2d2a2e"), "foreground": hexToRGBA("#f8f8f2"),
		"comment": hexToRGBA("#62606d"), "red": hexToRGBA("#ff657a"),
		"orange": hexToRGBA("#ffb270"), "yellow": hexToRGBA("#ffd76d"),
		"green": hexToRGBA("#bad761"), "cyan": hexToRGBA("#9aedfe"),
		"blue": hexToRGBA("#9aedfe"), "purple": hexToRGBA("#c39ac9"),
	},

	"machine": {
		"background": hexToRGBA("#2d2a2e"), "foreground": hexToRGBA("#f8f8f2"),
		"comment": hexToRGBA("#5c6370"), "red": hexToRGBA("#ff6e6e"),
		"orange": hexToRGBA("#f78c6c"), "yellow": hexToRGBA("#ffcb6b"),
		"green": hexToRGBA("#a9dc76"), "cyan": hexToRGBA("#78dce8"),
		"blue": hexToRGBA("#78dce8"), "purple": hexToRGBA("#ab9df2"),
	},

	"ristretto": {
		"background": hexToRGBA("#2c2525"), "foreground": hexToRGBA("#f2f2f2"),
		"comment": hexToRGBA("#6c6161"), "red": hexToRGBA("#ff657a"),
		"orange": hexToRGBA("#f79a32"), "yellow": hexToRGBA("#e6db74"),
		"green": hexToRGBA("#a6e22e"), "cyan": hexToRGBA("#66d9ef"),
		"blue": hexToRGBA("#66d9ef"), "purple": hexToRGBA("#ae81ff"),
	},
}
