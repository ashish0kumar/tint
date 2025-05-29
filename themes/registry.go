package themes

import (
	"fmt"
	"image/color"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var logInitOnce sync.Once

func setLogFlags() {
	log.SetFlags(0) // Simpler logging output
}

// hexToRGBA converts a hexadecimal color string ("#RRGGBB" or "#RRGGBBAA") to a color.RGBA
func hexToRGBA(hex string) color.RGBA {
	logInitOnce.Do(setLogFlags)

	if len(hex) < 7 || hex[0] != '#' {
		log.Fatalf("invalid hex color string: %s", hex)
	}

	hex = hex[1:] // Remove '#'

	if len(hex) != 6 && len(hex) != 8 {
		log.Fatalf("invalid hex color string length: %s", hex)
	}

	parseChannel := func(segment, name string) uint8 {
		val, err := parseHexChannel(segment)
		if err != nil {
			log.Fatalf("invalid hex color (%s channel): %s", name, hex)
		}
		return val
	}

	r := parseChannel(hex[0:2], "red")
	g := parseChannel(hex[2:4], "green")
	b := parseChannel(hex[4:6], "blue")
	a := uint8(0xFF)
	if len(hex) == 8 {
		a = parseChannel(hex[6:8], "alpha")
	}

	return color.RGBA{R: r, G: g, B: b, A: a}
}

// parseHexChannel parses a two-character hexadecimal string into a uint8 value
func parseHexChannel(s string) (uint8, error) {
	val, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		return 0, err
	}
	return uint8(val), nil
}

// validatePalette checks that a palette map contains valid colors
func validatePalette(paletteName string, palette map[string]color.RGBA) error {
	if len(palette) == 0 {
		return fmt.Errorf("palette '%s' is empty", paletteName)
	}
	if len(palette) < 3 {
		return fmt.Errorf("palette '%s' has too few colors (%d), need at least 3", paletteName, len(palette))
	}
	if len(palette) > 256 {
		return fmt.Errorf("palette '%s' has too many colors (%d), maximum is 256", paletteName, len(palette))
	}
	return nil
}

// AllThemeData holds all available themes
// The key is the theme name and value is a map of flavors to their color palettes
var AllThemeData = map[string]map[string]map[string]color.RGBA{
	"catppuccin": Catppuccin,
	"rosepine":   RosePine,
	"nord":       Nord,
	"tokyonight": TokyoNight,
	"gruvbox":    Gruvbox,
	"everforest": Everforest,
	"dracula":    Dracula,
	"solarized":  Solarized,
	"monochrome": Monochrome,
	"kanagawa":   Kanagawa,
	"ayu":        Ayu,
	"monokaipro": MonokaiPro,
	"nightowl":   NightOwl,
}

// GetPalette retrieves a palette by theme name and optional flavor
// Format: "theme-flavor" (e.g., "catppuccin-mocha")
func GetPalette(themeAndFlavor string) ([]color.Color, error) {
	cleaned := strings.ToLower(strings.TrimSpace(themeAndFlavor))
	if cleaned == "" {
		return nil, fmt.Errorf("theme name cannot be empty")
	}

	parts := strings.SplitN(cleaned, "-", 2)
	themeName := parts[0]
	subFlavor := ""
	if len(parts) > 1 {
		subFlavor = parts[1]
	}

	themeMap, ok := AllThemeData[themeName]
	if !ok {
		return nil, fmt.Errorf("invalid theme '%s'. Available themes: %s",
		themeName, strings.Join(GetAvailableThemeNames(), ", "))
	}

	var selectedPaletteMap map[string]color.RGBA
	if subFlavor != "" {
		if subPalette, ok := themeMap[subFlavor]; ok {
			selectedPaletteMap = subPalette
		} else {
			availableFlavors := GetAvailableFlavorNames(themeName)
			if len(availableFlavors) == 0 {
				return nil, fmt.Errorf("theme '%s' does not have flavors, use just '%s'", themeName, themeName)
			}
			return nil, fmt.Errorf("invalid flavor '%s' for theme '%s'. Available flavors: %s",
			subFlavor, themeName, strings.Join(availableFlavors, ", "))
		}
	} else {
		if defaultPalette, ok := themeMap["default"]; ok {
			selectedPaletteMap = defaultPalette
		} else {
			return nil, fmt.Errorf("theme '%s' has no defined palettes", themeName)
		}
	}

	paletteKey := themeName
	if subFlavor != "" {
		paletteKey = fmt.Sprintf("%s-%s", themeName, subFlavor)
	}

	if err := validatePalette(paletteKey, selectedPaletteMap); err != nil {
		return nil, fmt.Errorf("invalid palette for %s: %v", paletteKey, err)
	}

	// Convert to []color.Color
	paletteColors := make([]color.Color, 0, len(selectedPaletteMap))
	for _, c := range selectedPaletteMap {
		paletteColors = append(paletteColors, c)
	}

	return paletteColors, nil
}

// GetAvailableThemeNames returns a sorted slice of available theme names
func GetAvailableThemeNames() []string {
	names := make([]string, 0, len(AllThemeData))
	for name := range AllThemeData {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetAvailableFlavorNames returns a sorted slice of available flavor names for a given theme
func GetAvailableFlavorNames(themeName string) []string {
	themeMap, ok := AllThemeData[strings.ToLower(themeName)]
	if !ok {
		return nil
	}

	names := make([]string, 0, len(themeMap))
	for name := range themeMap {
		if name != "default" {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

// ValidateThemeData checks all theme data at startup
func ValidateThemeData() error {
	for themeName, themeMap := range AllThemeData {
		if len(themeMap) == 0 {
			return fmt.Errorf("theme '%s' has no flavor definitions", themeName)
		}

		for flavorName, palette := range themeMap {
			paletteKey := themeName
			if flavorName != "default" {
				paletteKey = fmt.Sprintf("%s-%s", themeName, flavorName)
			}

			if err := validatePalette(paletteKey, palette); err != nil {
				return err
			}
		}
	}
	return nil
}
