package themes

import (
	"fmt"
	"image/color"
	"log"
	"sort"
	"strconv"
	"strings"
)

// hexToRGBA converts a hexadecimal color string ("#RRGGBB") to a color.RGBA object
func hexToRGBA(hex string) color.RGBA {
	log.SetFlags(0)

	if len(hex) < 7 || hex[0] != '#' {
		log.Fatalf("invalid hex color string: %s", hex)
	}

	hex = hex[1:] // Remove '#'

	var r, g, b, a uint8
	var err error

	if len(hex) == 6 {

		// RGB format
		r, err = parseHexChannel(hex[0:2])
		if err != nil {
			log.Fatalf("invalid hex color (red channel): %s", hex)
		}
		g, err = parseHexChannel(hex[2:4])
		if err != nil {
			log.Fatalf("invalid hex color (green channel): %s", hex)
		}
		b, err = parseHexChannel(hex[4:6])
		if err != nil {
			log.Fatalf("invalid hex color (blue channel): %s", hex)
		}
		a = 0xFF // Default to fully opaque

	} else if len(hex) == 8 {

		// RGBA format
		r, err = parseHexChannel(hex[0:2])
		if err != nil {
			log.Fatalf("invalid hex color (red channel): %s", hex)
		}
		g, err = parseHexChannel(hex[2:4])
		if err != nil {
			log.Fatalf("invalid hex color (green channel): %s", hex)
		}
		b, err = parseHexChannel(hex[4:6])
		if err != nil {
			log.Fatalf("invalid hex color (blue channel): %s", hex)
		}
		a, err = parseHexChannel(hex[6:8])
		if err != nil {
			log.Fatalf("invalid hex color (alpha channel): %s", hex)
		}

	} else {
		log.Fatalf("invalid hex color string length: %s", hex)
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

	// Check for reasonable number of colors
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
	"ayu":        Ayu,
	"monokaipro": MonokaiPro,
	"nightowl":   NightOwl,
}

// GetPalette retrieves a palette by theme name and optional flavor
// It expects the format "theme-flavor" ("catppuccin-mocha")
func GetPalette(themeAndFlavor string) ([]color.Color, error) {
	if strings.TrimSpace(themeAndFlavor) == "" {
		return nil, fmt.Errorf("theme name cannot be empty")
	}

	parts := strings.SplitN(strings.ToLower(strings.TrimSpace(themeAndFlavor)), "-", 2)
	themeName := parts[0]
	subFlavor := ""
	if len(parts) > 1 {
		subFlavor = parts[1]
	}

	// Validate theme name
	if themeName == "" {
		return nil, fmt.Errorf("theme name cannot be empty")
	}

	themeMap, ok := AllThemeData[themeName]
	if !ok {
		availableThemes := GetAvailableThemeNames()
		return nil, fmt.Errorf("invalid theme '%s'. Available themes: %s",
			themeName, strings.Join(availableThemes, ", "))
	}

	var selectedPaletteMap map[string]color.RGBA

	if subFlavor != "" {
		// Specific flavor requested
		if subPalette, subOk := themeMap[subFlavor]; subOk {
			selectedPaletteMap = subPalette
		} else {
			availableFlavors := GetAvailableFlavorNames(themeName)
			if len(availableFlavors) == 0 {
				return nil, fmt.Errorf("theme '%s' does not have flavors, use just '%s'",
					themeName, themeName)
			}
			return nil, fmt.Errorf("invalid flavor '%s' for theme '%s'. Available flavors: %s",
				subFlavor, themeName, strings.Join(availableFlavors, ", "))
		}
	} else {
		// No flavor specified, find "default"
		if defaultPalette, defOk := themeMap["default"]; defOk {
			selectedPaletteMap = defaultPalette
		} else {
			return nil, fmt.Errorf("theme '%s' has no defined palettes", themeName)
		}
	}

	// Validate the selected palette
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

// GetAvailableFlavorNames returns a sorted slice of available flavor names for a given themeName
func GetAvailableFlavorNames(themeName string) []string {
	themeMap, ok := AllThemeData[strings.ToLower(themeName)]
	if !ok {
		return nil
	}

	names := make([]string, 0, len(themeMap))
	for name := range themeMap {
		if name == "default" { // Exclude "default"
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ValidateThemeData performs validation on all theme data at startup
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
