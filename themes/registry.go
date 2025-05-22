package themes

import (
	"fmt"
	"image/color"
	"sort"
	"strings"
)

// hexToRGBA converts a hex color string (like "#RRGGBB") to a color.RGBA struct
func hexToRGBA(hex string) color.RGBA {
	var r, g, b uint8
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		panic(fmt.Sprintf("Invalid hex color format: %s", hex))
	}
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		panic(fmt.Sprintf("Error parsing hex color %s: %v", hex, err))
	}
	return color.RGBA{R: r, G: g, B: b, A: 255} // Always opaque
}

// AllThemeData holds all available themes.
// The key is the theme name
// The value is a map of sub-flavors/variants to their color palettes.
var AllThemeData = map[string]map[string]map[string]color.RGBA{
	"catppuccin": CatppuccinPalettes,
	"rosepine":   RosePinePalettes,
	"nord":       NordPalettes,
	"tokyonight": TokyoNightPalettes,
	"gruvbox":    GruvboxPalettes,
	"everforest": EverforestPalettes,
	"dracula":    DraculaPalettes,
	"solarized":  SolarizedPalettes,
	"monochrome": MonochromePalettes,
}

// GetPalette retrieves a palette by theme name and optional sub-flavor
// It expects the format "theme-subflavor" (e.g., "catppuccin-mocha")
func GetPalette(themeAndFlavor string) ([]color.Color, error) {
	parts := strings.SplitN(strings.ToLower(themeAndFlavor), "-", 2)
	themeName := parts[0]
	subFlavor := ""
	if len(parts) > 1 {
		subFlavor = parts[1]
	}

	themeMap, ok := AllThemeData[themeName]
	if !ok {
		return nil, fmt.Errorf("invalid theme '%s'. Available themes: %s", themeName, strings.Join(GetAvailableThemeNames(), ", "))
	}

	var selectedPaletteMap map[string]color.RGBA
	if subFlavor != "" {
		if subPalette, subOk := themeMap[subFlavor]; subOk {
			selectedPaletteMap = subPalette
		} else {
			return nil, fmt.Errorf("invalid sub-flavor '%s' for theme '%s'. Available flavors: %s", subFlavor, themeName, strings.Join(GetAvailableSubFlavorNames(themeName), ", "))
		}
	} else {
		// If no sub-flavor specified, try to find a "default" or use the first available
		if defaultPalette, defOk := themeMap["default"]; defOk {
			selectedPaletteMap = defaultPalette
		} else if len(themeMap) > 0 {
			// If no "default", just pick the first available sub-palette (e.g., "mocha" for catppuccin if not specified)
			for _, palette := range themeMap {
				selectedPaletteMap = palette
				break // Take the first one
			}
		} else {
			return nil, fmt.Errorf("theme '%s' has no defined palettes", themeName)
		}
	}

	if len(selectedPaletteMap) == 0 {
		return nil, fmt.Errorf("no colors found for the selected theme/flavor: %s", themeAndFlavor)
	}

	paletteColors := make([]color.Color, 0, len(selectedPaletteMap))
	for _, c := range selectedPaletteMap {
		paletteColors = append(paletteColors, c)
	}
	return paletteColors, nil
}

// GetAvailableThemeNames returns a sorted slice of available top-level theme names.
func GetAvailableThemeNames() []string {
	names := make([]string, 0, len(AllThemeData))
	for name := range AllThemeData {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetAvailableSubFlavorNames returns a sorted slice of available sub-flavor names for a given themeName.
func GetAvailableSubFlavorNames(themeName string) []string {
	themeMap, ok := AllThemeData[strings.ToLower(themeName)]
	if !ok {
		return nil // Theme not found
	}

	names := make([]string, 0, len(themeMap))
	for name := range themeMap {
		if name == "default" { // Exclude "default" from the explicit flavor list
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
