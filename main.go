package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"tint/themes"
)

// toRGBA converts any color.Color to an 8-bit per channel color.RGBA
func toRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

// colorDistanceSquared calculates the squared Euclidean distance between two colors in RGB space
// This is faster than colorDistance as it avoids math.Sqrt and is sufficient for comparisons.
func colorDistanceSquared(c1, c2 color.RGBA) float64 {
	dr := float64(c1.R) - float64(c2.R)
	dg := float64(c1.G) - float64(c2.G)
	db := float64(c1.B) - float64(c2.B)

	return dr*dr + dg*dg + db*db
}

// findNClosestColors finds the N closest colors in the given palette to the original color
// It returns a slice of structs containing the distance and the color, sorted by distance
func findNClosestColors(originalRGBA color.RGBA, paletteRGBAs []color.RGBA, n int) []struct {
	dist  float64
	color color.Color // Keep color.Color for blending, as original input to blendColors is color.Color
} {
	if len(paletteRGBAs) == 0 {
		return nil
	}

	distances := make([]struct {
		dist  float64
		color color.Color // Store color.Color for consistency with blendColors
	}, len(paletteRGBAs))

	for i, pRGBA := range paletteRGBAs {
		distances[i] = struct {
			dist  float64
			color color.Color
		}{dist: colorDistanceSquared(originalRGBA, pRGBA), color: pRGBA} // Use squared distance here
	}

	sort.Slice(distances, func(i, j int) bool {
		return distances[i].dist < distances[j].dist
	})

	if n > len(distances) {
		n = len(distances)
	}
	return distances[:n]
}

// blendColors takes a slice of colors and their corresponding weights and returns a single blended color
func blendColors(colors []color.Color, weights []float64) color.RGBA {
	if len(colors) == 0 || len(colors) != len(weights) {
		return color.RGBA{}
	}

	var sumR, sumG, sumB float64
	var totalWeight float64

	for i := range colors {
		rgba := toRGBA(colors[i])
		sumR += float64(rgba.R) * weights[i]
		sumG += float64(rgba.G) * weights[i]
		sumB += float64(rgba.B) * weights[i]
		totalWeight += weights[i]
	}

	if totalWeight == 0 {
		return toRGBA(colors[0]) // Fallback to the first color if weights are somehow zero
	}

	return color.RGBA{
		R: uint8(math.Round(sumR / totalWeight)),
		G: uint8(math.Round(sumG / totalWeight)),
		B: uint8(math.Round(sumB / totalWeight)),
		A: 255,
	}
}

// applyLuminosity adjusts a color's brightness by scaling its RGB components
func applyLuminosity(c color.RGBA, factor float64) color.RGBA {
	r := uint8(math.Max(0, math.Min(255, float64(c.R)*factor)))
	g := uint8(math.Max(0, math.Min(255, float64(c.G)*factor)))
	b := uint8(math.Max(0, math.Min(255, float64(c.B)*factor)))
	return color.RGBA{R: r, G: g, B: b, A: c.A}
}

// shepardsMethodColor applies Shepard's Method for color interpolation
// It finds the 'nearest' palette colors and blends them using inverse distance weighting
func shepardsMethodColor(originalRGBA color.RGBA, paletteRGBAs []color.RGBA, nearest int, power float64) color.Color {
	closest := findNClosestColors(originalRGBA, paletteRGBAs, nearest)
	if len(closest) == 0 {
		return originalRGBA // Return original RGBA as a Color
	}
	// If an exact match is found or only one neighbor is requested, just return it
	if len(closest) == 1 || closest[0].dist == 0 {
		return closest[0].color
	}

	weights := make([]float64, len(closest))
	var totalWeight float64
	for i, c := range closest {
		if c.dist == 0 { // Avoid division by zero if there's an exact match
			return c.color
		}
		// Weight is inversely proportional to distance raised to the power
		// Use the actual distance (sqrt of squared distance) for inverse distance weighting
		weight := 1.0 / math.Pow(math.Sqrt(c.dist), power)
		weights[i] = weight
		totalWeight += weight
	}

	if totalWeight == 0 { // Fallback if all weights somehow sum to zero
		return closest[0].color
	}

	return blendColors(extractColors(closest), weights)
}

// extractColors pulls just the color.Color from the sorted slice of (distance, color) tuples
func extractColors(sortedColors []struct {
	dist  float64
	color color.Color
}) []color.Color {
	colors := make([]color.Color, len(sortedColors))
	for i, item := range sortedColors {
		colors[i] = item.color
	}
	return colors
}

// processImageWithShepardsMethod applies Shepard's Method to each pixel of the image concurrently
func processImageWithShepardsMethod(img image.Image, palette []color.Color, luminosity float64, nearest int, power float64) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	// Pre-convert the palette to RGBA once
	paletteRGBAs := make([]color.RGBA, len(palette))
	for i, c := range palette {
		paletteRGBAs[i] = toRGBA(c)
	}

	// Determine the number of goroutines to use
	numWorkers := runtime.GOMAXPROCS(0) // Use number of logical CPUs
	if numWorkers == 0 {
		numWorkers = 1 // Fallback if for some reason GOMAXPROCS returns 0
	}
	if numWorkers > bounds.Dy() {
		numWorkers = bounds.Dy() // Don't create more workers than rows
	}

	// Divide the image into horizontal chunks for workers
	rowsPerWorker := bounds.Dy() / numWorkers
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			startY := bounds.Min.Y + workerID*rowsPerWorker
			endY := startY + rowsPerWorker
			if workerID == numWorkers-1 { // Last worker takes remaining rows
				endY = bounds.Max.Y
			}

			for y := startY; y < endY; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					originalColor := img.At(x, y)
					originalRGBA := toRGBA(originalColor) // Convert to RGBA once per pixel
					a := originalRGBA.A

					if a == 0 { // If the pixel is fully transparent, keep it that way
						newImg.Set(x, y, color.Transparent)
					} else {
						// Adjust luminosity first
						adjustedColor := applyLuminosity(originalRGBA, luminosity)
						finalColor := shepardsMethodColor(adjustedColor, paletteRGBAs, nearest, power)
						newImg.Set(x, y, finalColor)
					}
				}
			}
		}(i)
	}

	wg.Wait() // Wait for all to complete
	return newImg
}

// listThemes prints all available themes and their flavors
func listThemes() {
	fmt.Println("Available Themes and Flavors:")

	// Get sorted list of theme names
	themeNames := themes.GetAvailableThemeNames()

	for _, themeName := range themeNames {
		fmt.Printf("  %s:\n", themeName)
		flavors := themes.GetAvailableSubFlavorNames(themeName)

		if len(flavors) == 0 {
			fmt.Println("    (No specific flavors, use as default)")
		} else {
			// Sort flavors for consistent output
			sort.Strings(flavors)
			for _, flavor := range flavors {
				// Indicate default flavor if it exists
				if flavor == "default" {
					fmt.Printf("    - %s (default)\n", themeName) // Ex: catppuccin (default)
				} else {
					fmt.Printf("    - %s-%s\n", themeName, flavor) // Ex: catppuccin-latte
				}
			}
		}
	}
	fmt.Println("\nUsage example: tint --theme <theme-name> or tint --theme <theme-name>-<flavor>")
}

func main() {
	log.SetFlags(0)

	// --- Define variables for flags ---

	var imagePath string
	var themeAndFlavor string
	var outputPath string
	var jpegQuality int
	var luminosity float64
	var nearest int
	var power float64
	var listThemesFlag bool

	// --- Define and Parse CLI Flags ---

	flag.StringVar(&imagePath, "image", "", "Path to the input image (required). Supports JPEG, PNG")
	flag.StringVar(&imagePath, "i", "", "Shorthand for -image")

	flag.StringVar(&themeAndFlavor, "theme", "catppuccin-mocha", "Theme palette and optional flavor. Use -list-themes or -l to see all options.")
	flag.StringVar(&themeAndFlavor, "t", "catppuccin-mocha", "Shorthand for -theme")

	flag.StringVar(&outputPath, "output", "", "Path for the output image (default: <input_filename>_themed_<theme-flavor>.png)")
	flag.StringVar(&outputPath, "o", "", "Shorthand for -output")

	flag.IntVar(&jpegQuality, "jpeg-quality", 80, "JPEG quality for output (1-100), only for JPEG output")

	// Parameters specific to Shepard's Method
	flag.Float64Var(&luminosity, "luminosity", 1.0, "Luminosity adjustment factor (e.g., 0.8 for darker, 1.2 for brighter)")
	flag.IntVar(&nearest, "nearest", 26, "Number of nearest palette colors to consider for interpolation")
	flag.Float64Var(&power, "power", 4.0, "Power for Shepard's Method (influences how quickly weights fall off)")

	flag.BoolVar(&listThemesFlag, "list-themes", false, "List all available themes and their flavors")
	flag.BoolVar(&listThemesFlag, "l", false, "Shorthand for -list-themes")

	flag.Parse()

	// --- Handle listThemesFlag ---

	if listThemesFlag {
		listThemes()
		os.Exit(0) // Exit after listing themes, no further processing needed
	}

	// --- Validate Inputs ---

	if imagePath == "" {
		log.Println("Usage: tint --image <path_to_image> [options]")
		log.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Get the palette using the themes package function
	paletteColors, err := themes.GetPalette(themeAndFlavor)
	if err != nil {
		log.Fatalf("Error getting palette: %v", err)
	}

	if jpegQuality < 1 || jpegQuality > 100 {
		log.Fatalf("Error: JPEG quality must be between 1 and 100. You entered: %d", jpegQuality)
	}

	// --- Open and Decode Image ---

	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Error opening image '%s': %v", imagePath, err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Error decoding image '%s': %v. Make sure it's a JPEG or PNG file.", imagePath, err)
	}
	if format != "jpeg" && format != "png" {
		log.Fatalf("Error: Unsupported input image format '%s'. We only support JPEG and PNG.", format)
	}
	log.Printf("Image '%s' loaded successfully. Format: %s", imagePath, format)

	// --- Process Image with Shepard's Method ---

	log.Printf("Processing image with Shepard's Method (theme: %s, nearest: %d, power: %.1f)", themeAndFlavor, nearest, power)
	processedImg := processImageWithShepardsMethod(img, paletteColors, luminosity, nearest, power)

	// --- Determine Output Path and Save Image ---

	outPath := outputPath
	if outPath == "" {
		dir := filepath.Dir(imagePath)
		base := filepath.Base(imagePath)
		nameWithoutExt := strings.TrimSuffix(base, filepath.Ext(base))
		outPath = filepath.Join(dir, fmt.Sprintf("%s_themed_%s%s", nameWithoutExt, strings.ToLower(themeAndFlavor), ".png"))
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("Error creating output file '%s': %v", outPath, err)
	}
	defer outFile.Close()

	outputExt := strings.ToLower(filepath.Ext(outPath))
	chosenOutputFormat := ""

	if outputExt == ".jpg" || outputExt == ".jpeg" {
		chosenOutputFormat = "jpeg"
	} else if outputExt == ".png" {
		chosenOutputFormat = "png"
	} else {
		log.Printf("Warning: Output path '%s' doesn't have a .png, .jpg, or .jpeg extension. Saving as PNG.", outPath)
		chosenOutputFormat = "png"
		if outputPath == "" {
			outPath = strings.TrimSuffix(outPath, filepath.Ext(outPath)) + ".png"
			outFile.Close()
			outFile, err = os.Create(outPath)
			if err != nil {
				log.Fatalf("Error re-creating output file with .png extension '%s': %v", outPath, err)
			}
			defer outFile.Close()
			log.Printf("Adjusted output path to: %s", outPath)
		}
	}

	log.Printf("Trying to save the image as %s to %s", chosenOutputFormat, outPath)

	switch chosenOutputFormat {
	case "jpeg":
		var opt jpeg.Options
		opt.Quality = jpegQuality
		if err := jpeg.Encode(outFile, processedImg, &opt); err != nil {
			log.Fatalf("Error encoding JPEG image to '%s': %v", outPath, err)
		}
	case "png":
		if err := png.Encode(outFile, processedImg); err != nil {
			log.Fatalf("Error encoding PNG image to '%s': %v", outPath, err)
		}
	default:
		log.Fatalf("Oops! Couldn't determine a supported output format: %s", chosenOutputFormat)
	}

	log.Printf("Image successfully themed and saved to '%s' with the '%s' theme", outPath, strings.ToLower(themeAndFlavor))
}
