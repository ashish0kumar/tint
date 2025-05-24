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
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ashish0kumar/tint/themes"
)

const (
	MaxImageDimension = 10000    // Maximum allowed width or height of the image
	MaxImagePixels    = 50000000 // Maximum allowed number of pixels in the image (7001x7001)

	// Default params for Shepard's Method
	defaultLuminosity = 1.0
	defaultNearest    = 30
	defaultPower      = 2.5

	// ANSI escape codes for formatting
	bold      = "\033[1m"
	underline = "\033[4m"
	reset     = "\033[0m"
)

// toRGBA converts any color.Color to color.RGBA
func toRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

// colorDistanceSquared calculates the squared euclidean distance between two colors in RGB space
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
	color color.Color
} {
	if len(paletteRGBAs) == 0 {
		return nil
	}

	distances := make([]struct {
		dist  float64
		color color.Color
	}, 0, len(paletteRGBAs))

	for _, pRGBA := range paletteRGBAs {
		distances = append(distances, struct {
			dist  float64
			color color.Color
		}{dist: colorDistanceSquared(originalRGBA, pRGBA), color: pRGBA})
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
		return originalRGBA // No palette colors available, return original color
	}
	// If an exact match is found or only one neighbor is requested, just return it
	if len(closest) == 1 || closest[0].dist == 0 {
		return closest[0].color
	}

	weights := make([]float64, len(closest))
	var totalWeight float64
	for i, c := range closest {
		if c.dist == 0 { // Avoid division by zero
			return c.color
		}
		// Inverse distance weighting
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

// ProgressTracker tracks and displays processing progress
type ProgressTracker struct {
	total       int64
	processed   int64
	startTime   time.Time
	lastUpdate  time.Time
	updateMutex sync.Mutex
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(total int64) *ProgressTracker {
	return &ProgressTracker{
		total:      total,
		startTime:  time.Now(),
		lastUpdate: time.Now(),
	}
}

// updateProgress increments the processed count and displays progress
func (pt *ProgressTracker) updateProgress(increment int64) {
	atomic.AddInt64(&pt.processed, increment)

	pt.updateMutex.Lock()
	defer pt.updateMutex.Unlock()

	now := time.Now()
	if now.Sub(pt.lastUpdate) < 100*time.Millisecond {
		return
	}
	pt.lastUpdate = now

	processed := atomic.LoadInt64(&pt.processed)
	if processed >= pt.total {
		return
	}

	percentage := float64(processed) / float64(pt.total) * 100
	elapsed := now.Sub(pt.startTime)

	if processed > 0 {
		estimatedTotal := time.Duration(float64(elapsed) / float64(processed) * float64(pt.total))
		remaining := estimatedTotal - elapsed

		fmt.Printf("\rProgress: %.1f%% (%d/%d) Elapsed: %v ETA: %v",
			percentage, processed, pt.total, elapsed.Round(time.Second), remaining.Round(time.Second))
	}
}

// finishProgress completes the progress display
func (pt *ProgressTracker) finishProgress() {
	processed := atomic.LoadInt64(&pt.processed)
	elapsed := time.Since(pt.startTime)
	fmt.Printf("\rComplete: 100.0%% (%d/%d) in %v\n",
		processed, pt.total, elapsed.Round(time.Millisecond))
}

// processImageWithShepardsMethod applies Shepard's Method to each pixel of the image concurrently
func processImageWithShepardsMethod(img image.Image, palette []color.Color, luminosity float64, nearest int, power float64) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	// Pre convert the palette to RGBA once
	paletteRGBAs := make([]color.RGBA, len(palette))
	for i, c := range palette {
		paletteRGBAs[i] = toRGBA(c)
	}

	// Initialize progress tracker
	totalPixels := int64(bounds.Dx() * bounds.Dy())
	progress := NewProgressTracker(totalPixels)

	// Implement concurrency
	numWorkers := runtime.GOMAXPROCS(0)
	if numWorkers == 0 {
		numWorkers = 1
	}
	if numWorkers > bounds.Dy() {
		numWorkers = bounds.Dy() // Dont create more workers than rows
	}

	// Divide the image into horizontal chunks for workers
	rowsPerWorker := bounds.Dy() / numWorkers
	var wg sync.WaitGroup

	for i := range numWorkers {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			startY := bounds.Min.Y + workerID*rowsPerWorker
			endY := startY + rowsPerWorker
			if workerID == numWorkers-1 { // Last worker takes remaining rows
				endY = bounds.Max.Y
			}

			pixelsProcessed := int64(0)
			for y := startY; y < endY; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					originalColor := img.At(x, y)
					originalRGBA := toRGBA(originalColor) // Convert to RGBA once per pixel
					a := originalRGBA.A

					if a == 0 { // If the pixel is fully transparent, keep it that way
						newImg.Set(x, y, color.Transparent)
					} else {
						// Adjust luminosity and apply Shepard's Method
						adjustedColor := applyLuminosity(originalRGBA, luminosity)
						finalColor := shepardsMethodColor(adjustedColor, paletteRGBAs, nearest, power)
						newImg.Set(x, y, finalColor)
					}
					pixelsProcessed++
				}
				// Update progress every row to avoid too frequent updates
				if y%10 == 0 || y == endY-1 {
					progress.updateProgress(pixelsProcessed)
					pixelsProcessed = 0
				}
			}
			// Update remaining pixels
			if pixelsProcessed > 0 {
				progress.updateProgress(pixelsProcessed)
			}
		}(i)
	}

	wg.Wait() // Wait for all to complete
	progress.finishProgress()
	return newImg
}

// validateInputs performs extensive input validation
func validateInputs(imagePath string, themeAndFlavor string, luminosity float64, nearest int, power float64) error {
	// Check if image file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("image file '%s' does not exist", imagePath)
	}

	// Check file permissions
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("cannot open image file '%s': %v", imagePath, err)
	}
	defer file.Close()

	// Get file info for size validation
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("cannot get file info for '%s': %v", imagePath, err)
	}

	// Check file size
	if fileInfo.Size() > 100*1024*1024 {
		return fmt.Errorf("image file '%s' is too large (%.2f MB). Maximum size is 100 MB",
			imagePath, float64(fileInfo.Size())/(1024*1024))
	}

	// Try to decode image to check format and dimensions
	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("cannot decode image '%s': %v. Make sure it's a valid JPEG or PNG file", imagePath, err)
	}
	if format != "jpeg" && format != "png" {
		return fmt.Errorf("unsupported image format '%s'. Only JPEG and PNG are supported", format)
	}

	// Check image dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	totalPixels := width * height

	if width > MaxImageDimension || height > MaxImageDimension {
		return fmt.Errorf("image dimensions too large (%dx%d). Maximum dimension is %d pixels",
			width, height, MaxImageDimension)
	}
	if totalPixels > MaxImagePixels {
		return fmt.Errorf("image has too many pixels (%d). Maximum is %d pixels",
			totalPixels, MaxImagePixels)
	}

	// Validate theme exists
	if _, err := themes.GetPalette(themeAndFlavor); err != nil {
		return fmt.Errorf("theme validation failed: %v", err)
	}

	// Validate params
	if luminosity <= 0 {
		return fmt.Errorf("luminosity must be positive, got %.2f", luminosity)
	}
	if nearest < 1 {
		return fmt.Errorf("nearest colors count must be at least 1, got %d", nearest)
	}
	if power <= 0 {
		return fmt.Errorf("power must be positive, got %.2f", power)
	}
	// log.Printf("Image validation passed: %dx%d pixels, %s format, %.2f MB",
	// width, height, strings.ToUpper(format), float64(fileInfo.Size())/(1024*1024))

	return nil
}

// loadImage loads and returns the image
func loadImage(imagePath string) (image.Image, string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, "", fmt.Errorf("error opening image '%s': %v", imagePath, err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", fmt.Errorf("error decoding image '%s': %v", imagePath, err)
	}

	return img, format, nil
}

// getOutputExtension determines the output file extension based on input format and output path
func getOutputExtension(inputFormat string, outputPath string) string {
	if outputPath != "" {
		// If output path is specified, use its extension
		ext := strings.ToLower(filepath.Ext(outputPath))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
			return ext
		}
	}

	// Preserve input format
	switch inputFormat {
	case "jpeg":
		return ".jpg"
	case "png":
		return ".png"
	default:
		return ".png" // Fallback
	}
}

// saveImage saves the processed image
func saveImage(img image.Image, outputPath string, inputFormat string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file '%s': %v", outputPath, err)
	}
	defer outFile.Close()

	outputExt := getOutputExtension(inputFormat, outputPath)

	switch outputExt {
	case ".jpg", ".jpeg":
		var opt jpeg.Options
		opt.Quality = 85
		if err := jpeg.Encode(outFile, img, &opt); err != nil {
			return fmt.Errorf("error encoding JPEG image to '%s': %v", outputPath, err)
		}
	case ".png":
		if err := png.Encode(outFile, img); err != nil {
			return fmt.Errorf("error encoding PNG image to '%s': %v", outputPath, err)
		}
	default:
		return fmt.Errorf("unsupported output format '%s'. Use .png, .jpg, or .jpeg", outputExt)
	}

	return nil
}

// generateOutputPath creates the output path based on input path, theme and format
func generateOutputPath(inputPath string, themeAndFlavor string, inputFormat string) string {
	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)
	nameWithoutExt := strings.TrimSuffix(base, filepath.Ext(base))

	ext := getOutputExtension(inputFormat, "")

	return filepath.Join(dir, fmt.Sprintf("%s_themed_%s%s",
		nameWithoutExt, strings.ToLower(themeAndFlavor), ext))
}

// listThemes prints all available themes and their flavors
func listThemes() {
	fmt.Printf("\n%s%sUsage:%s tint -i <IMAGE> -t <THEME-FLAVOR> [OPTIONS]\n\n", bold, underline, reset)
	fmt.Printf("%s%sAvailable Themes & Flavors:%s\n", bold, underline, reset)

	themeNames := themes.GetAvailableThemeNames()

	for _, themeName := range themeNames {
		flavors := themes.GetAvailableFlavorNames(themeName)
		fmt.Printf("\n  %s%s%s\n", bold, themeName, reset)

		if len(flavors) > 0 {
			for _, flavor := range flavors {
				fmt.Printf("    - %s-%s\n", themeName, flavor)
			}
		}
	}

	fmt.Println()
}

// openFileInDefaultViewer attempts to open a file using the OS default viewer.
func openFileInDefaultViewer(filePath string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", filePath)
	case "windows": // Windows
		cmd = exec.Command("cmd", "/c", "start", filePath)
	case "linux": // Linux
		cmd = exec.Command("xdg-open", filePath)
	default:
		log.Printf("Unsupported operating system for automatic file opening: %s", runtime.GOOS)
		return
	}

	err := cmd.Start()
	if err != nil {
		log.Printf("Failed to open file '%s' in default viewer: %v", filePath, err)
	} else {
		log.Printf("Opened '%s' in default viewer.", filePath)
	}
}

func main() {
	log.SetFlags(0)
	themes.ValidateThemeData() // Validate theme data at startup

	// --- Define variables for flags ---

	var imagePath string
	var themeAndFlavor string
	var outputPath string
	var luminosity float64
	var nearest int
	var power float64
	var listThemesFlag bool

	// --- Define and parse flags ---

	flag.StringVar(&imagePath, "image", "", "Path to the input image (required). Supports JPEG, PNG")
	flag.StringVar(&imagePath, "i", "", "Shorthand for -image")

	flag.StringVar(&themeAndFlavor, "theme", "", "Theme palette and optional flavor. Use -list-themes or -l to see all options.")
	flag.StringVar(&themeAndFlavor, "t", "", "Shorthand for -theme")

	flag.StringVar(&outputPath, "output", "", "Path for the output image (default: <input_filename>_themed_<theme-flavor>.<input_format>)")
	flag.StringVar(&outputPath, "o", "", "Shorthand for -output")

	flag.BoolVar(&listThemesFlag, "list-themes", false, "List all available themes and their flavors")
	flag.BoolVar(&listThemesFlag, "l", false, "Shorthand for -list-themes")

	// Params specific to Shepard's Method
	flag.Float64Var(&luminosity, "luminosity", defaultLuminosity, "Luminosity adjustment factor (e.g., 0.8 for darker, 1.2 for brighter)")
	flag.IntVar(&nearest, "nearest", defaultNearest, "Number of nearest palette colors to consider for interpolation")
	flag.Float64Var(&power, "power", defaultPower, "Power for Shepard's Method (influences how quickly weights fall off)")

	flag.Usage = setUsage

	flag.Parse()

	// --- Handle listThemesFlag ---

	if listThemesFlag {
		listThemes()
		os.Exit(0)
	}

	// --- Validate required args ---

	if imagePath == "" {
		fmt.Fprintln(os.Stderr, "Error: -i or --image <IMAGE> is required.")
		flag.Usage()
		os.Exit(1)
	}

	if themeAndFlavor == "" {
		fmt.Fprintln(os.Stderr, "Error: -t or --theme <THEME-FLAVOR> is required.")
		flag.Usage()
		os.Exit(1)
	}

	// --- Input validation ---

	// log.Println("Validating inputs...")
	if err := validateInputs(imagePath, themeAndFlavor, luminosity, nearest, power); err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	// --- Get palette ---

	paletteColors, err := themes.GetPalette(themeAndFlavor)
	if err != nil {
		log.Fatalf("Error getting palette: %v", err)
	}

	// --- Load image ---

	// log.Printf("Loading image '%s'...", imagePath)
	img, format, err := loadImage(imagePath)
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}
	// log.Printf("Image loaded successfully. Format: %s", format)

	// --- Process image with shepard's method ---

	// log.Printf("Processing image with Shepard's Method (theme: %s, nearest: %d, power: %.1f)", themeAndFlavor, nearest, power)

	log.Printf("Theme: %s", strings.ToLower(themeAndFlavor))
	log.Printf("Shepard's Method: nearest = %d, power = %.1f, luminosity = %.1f", nearest, power, luminosity)
	log.Printf("Processing: '%s'", imagePath)

	processedImg := processImageWithShepardsMethod(img, paletteColors, luminosity, nearest, power)

	// --- Determine output path ---

	outPath := outputPath
	if outPath == "" {
		outPath = generateOutputPath(imagePath, themeAndFlavor, format)
	}

	// --- Save image ---

	// log.Printf("Saving processed image to '%s'...", outPath)
	if err := saveImage(processedImg, outPath, format); err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}

	log.Printf("Saved image: '%s'\n", outPath)

	// --- Open output image in default viewer ---
	openFileInDefaultViewer(outPath)
}

// setUsage prints the custom help message
func setUsage() {
	programName := filepath.Base(os.Args[0])
	w := flag.CommandLine.Output()

	// Usage
	fmt.Fprintf(w, "\n%s%sUsage:%s %s -i <IMAGE> -t <THEME-FLAVOR> [OPTIONS]\n\n", bold, underline, reset, programName)

	// Theme
	fmt.Fprintf(w, "  %s-t, --theme <STRING>%s\n", bold, reset)
	fmt.Fprintf(w, "\tTheme palette and optional flavor (required).\n")
	fmt.Fprintf(w, "\tUse -l or --list-themes to see all available themes and flavors.\n\n")

	// Image
	fmt.Fprintf(w, "  %s-i, --image <PATH>%s\n", bold, reset)
	fmt.Fprintf(w, "\tPath to the input image (required). Supports JPEG, PNG.\n\n")

	// Options heading
	fmt.Fprintf(w, "%s%sOptions:%s\n\n", bold, underline, reset)

	// Output
	fmt.Fprintf(w, "  %s-o, --output <PATH>%s\n", bold, reset)
	fmt.Fprintf(w, "\tPath for the output image.\n")
	fmt.Fprintf(w, "\t(Default: <input_filename>_themed_<theme-flavor>.<input_format>)\n\n")

	// Luminosity
	fmt.Fprintf(w, "  %s--luminosity <FLOAT>%s\n", bold, reset)
	fmt.Fprintf(w, "\tLuminosity adjustment factor (e.g., 0.8 for darker, 1.2 for brighter).\n")
	fmt.Fprintf(w, "\t(Default: %.1f)\n\n", defaultLuminosity)

	// Nearest
	fmt.Fprintf(w, "  %s--nearest <COUNT>%s\n", bold, reset)
	fmt.Fprintf(w, "\tNumber of nearest palette colors to consider for interpolation.\n")
	fmt.Fprintf(w, "\t(Default: %d)\n\n", defaultNearest)

	// Power
	fmt.Fprintf(w, "  %s--power <FLOAT>%s\n", bold, reset)
	fmt.Fprintf(w, "\tPower for Shepard's Method (influences how quickly weights fall off).\n")
	fmt.Fprintf(w, "\t(Default: %.1f)\n\n", defaultPower)

	// List Themes
	fmt.Fprintf(w, "  %s-l, --list-themes%s\n", bold, reset)
	fmt.Fprintf(w, "\tList all available themes and their flavors.\n\n")

	// Help
	fmt.Fprintf(w, "  %s-h, --help%s\n", bold, reset)
	fmt.Fprintf(w, "\tPrint this help message.\n")
}
