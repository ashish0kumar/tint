# tint

`tint` is a command-line tool to **recolor images using your favorite theme palettes**. It's designed for those who appreciate a cohesive visual aesthetic, letting you match your wallpapers or other images to your favorite themes.

## Features

- **Theme-Based Recolorization:** Apply color palettes from themes like Catppuccin, Nord, Gruvbox, and many more.
- **Smooth Color Transitions:** Uses Shepard's Method for natural gradients and blends in complex images.
- **Luminosity Adjustment:** Easily fine-tune the brightness of your recolored images.
- **Customizable Interpolation:** Control blending by adjusting `nearest` colors and weighting function `power`.
- **Image Format Support:** Works with JPEG and PNG image files.
- **Efficient Processing:**  Leverages Go's concurrency for quick processing, specially for large images.
- **Lightweight & Dependency-Free:** A single, self-contained Go binary with no external dependencies.

## Available Themes

- Catppuccin
- Dracula
- Everforest
- Gruvbox
- Monochrome
- Nord
- RosePine
- Solarized
- Tokyonight

## Installation

### Prerequisites:

[Go 1.18+](https://golang.org/doc/install)

### Install via `go install`

```bash
go install github.com/ashish0kumar/tint@latest
```

### Build from Source

Clone the repo, build the project and move it inside your `$PATH`

```Bash
git clone https://github.com/ashish0kumar/tint.git
cd tint
go build
sudo mv tint /usr/local/bin/
```

## Usage

```Bash
tint -i <IMAGE_PATH> -t <THEME-FLAVOR> [OPTIONS]
```

### Required Arguments

- `-i, --image <PATH>`: Specifies the path to your input image, supports JPEG and PNG formats.
- `-t, --theme <STRING>`: Defines the theme palette and its optional flavor to apply. Use `-l` or `--list-themes` to see all available options.

### Options

- `-o, --output <PATH>`: Sets the path for the output image.
    - (Default: `<input_filename>_themed_<theme-flavor>.<input_format>`)

- `--luminosity <FLOAT>`: Adjusts the overall luminosity (brightness) of the image. A value of `0.8` makes it darker, `1.2` makes it brighter.
    - (Default: `1.0`)

- `--nearest <COUNT>`: Determines the number of nearest palette colors considered for interpolation. A higher count can result in smoother blending.
    - (Default: `26`)

- `--power <FLOAT>`: This parameter for Shepard's Method influences how quickly the weights of distant colors fall off. Higher values mean closer colors have a stronger influence.
    - (Default: `4.0`)

- `-l, --list-themes`: Displays all available themes and their flavors.

- `-h, --help`: Shows the command-line help message.

## Examples

- **List all available themes and flavors:**

```Bash
tint -l
```

- **Recolor an image with the Nord theme (default flavor):**

```Bash
tint -i my_wallpaper.jpg -t nord
```

- **Apply the Catppuccin Mocha theme and save to a specific path:**

```Bash
tint -i original.png -t catppuccin-mocha -o catppuccin_output.jpg
```

- **Recolor an image with Gruvbox Dark, making it slightly brighter:**

```Bash
tint -i photo.jpeg -t gruvbox-dark --luminosity 1.2
```

- **Use fewer nearest colors for a potentially more distinct mapping:**

```Bash
tint -i gradient_art.png -t rosepine --nearest 10
```

## Development

### Submitting New Themes

I'd love to expand `tint`'s theme collection! If you have a favorite theme not yet included, or want to contribute a new one, here's how:

1. **Understand the structure:**
    - Theme definitions live in the `themes/` directory. Each theme typically gets its own `.go` file (eg `themes/catppuccin.go`).

    - Colors are defined as `color.RGBA` values, usually converted from hexadecimal strings using the `hexToRGBA` helper function found in `themes/registry.go`

2. **Create your theme file:**

    - Create a new Go file in the `themes/` directory (eg `themes/mytheme.go`).

    - Inside this file, define a `map[string]map[string]color.RGBA` that holds your theme data.

        - The top-level key should be your theme's main name (eg `"mytheme"`).
        
        - The nested map contains flavors. If your theme has multiple flavors (eg dark/light variants), define them as sub-maps (eg `"dark": {...}, "light": {...}`).

        - **Crucially, every theme must include a `"default"` flavor**. This is the palette tint will use if no specific flavor is mentioned by the user when running the command (eg just `-t catppuccin` will pick `catppuccin-default`).
        
        - Map descriptive color names (eg `"base"`, `"surface0"`) to their `color.RGBA` values using `hexToRGBA`

    - **Example `themes/mytheme.go` structure:**
    
    ```Go
    package themes

    import "image/color"

    var MyTheme = map[string]map[string]color.RGBA{
        "default": { // Important
            "background": hexToRGBA("#282A36"),
            "foreground": hexToRGBA("#F8F8F2"),
            "comment":    hexToRGBA("#6272A4"),
            // ... more colors
        },
        "dark": {
            "base":     hexToRGBA("#1A1A1A"),
            "text":     hexToRGBA("#F0F0F0"),
            "accent":   hexToRGBA("#FF5733"),
            // ... more colors
        },
        "light": {
            "base":     hexToRGBA("#F0F0F0"),
            "text":     hexToRGBA("#1A1A1A"),
            "accent":   hexToRGBA("#337AFF"),
            // ... more colors
        },
    }
    ```

3. **Register your theme:**

    - Open `themes/registry.go`
    - Add your theme to the `AllThemeData` map.


    ```Go
    var AllThemeData = map[string]map[string]map[string]color.RGBA{
        // ... existing themes
        "mytheme": MyTheme, // Add this line
    }
    ```

4. **Validate and Test:**
    - Run `go build` from the project root to ensure there are no compilation errors.
    - Test your new theme using `tint -t mytheme-dark` (or `mytheme`) with an image to confirm it works as expected.

5. **Submit a Pull Request:**

    - Fork the repository.
    - Create a new branch for your changes.
    - Commit your new theme file and the changes to `themes/registry.go`
    - Open a Pull Request, explaining your new theme.

## Deep Dive

### What Makes `tint` Different

There are many excellent image recoloring tools out there, and each has its unique strengths and approaches. For instance, tools like [`dipc`](https://github.com/doprz/dipc) and [`faerber`](https://github.com/nekowinston/faerber) often use direct color mapping strategies (often based on *DeltaE* color differences). These methods are very **effective for images with consistent, fewer colors**, providing quick and precise color reassignments to match a target palette.

However, when an image features **complex gradients, subtle color mixtures, or a broad spectrum of colors**, direct mapping can sometimes result in a slightly "**patchy**" or less fluid appearance. This is where `tint` offers a distinct approach.

`tint` employs **Shepard's Method (Inverse Distance Weighting)** for its color interpolation. Instead of simply replacing a pixel's color with the single closest match from the palette, `tint` considers **multiple nearest palette colors** and creates a **weighted average**. The closer a palette color is to the original, the more influence it has on the final pixel color. This blending technique provides:

- **Smoother gradients:** Complex color transitions in your original image are gracefully transformed into the new palette without harsh banding.

- **Natural appearance:** The resulting images maintain a more organic and less "digital" feel, making them well-suited for a variety of applications where the fidelity of smooth areas is important.

Beyond its unique interpolation method, `tint` is also specifically designed to be **incredibly lightweight and self-contained**. Built purely with Go's standard library, it has **no external dependencies**. This focus on minimalism and efficiency complements other valuable recoloring utilities that might offer broader feature sets but come with more overhead.


### How It Works

`tint` implements **Shepard's Method (Inverse Distance Weighting)** for its core color transformation. Here's a simplified breakdown:

- **Initial luminosity adjustment:** If specified, the pixel's color components are scaled to adjust its overall brightness.

- **Nearest neighbor search:** For each pixel in the input image, `tint` identifies a specified number (`--nearest`) of the closest colors from the chosen theme's palette, based on their RGB distance.

- **Weighted contribution:** Each identified palette color is assigned a weight. This weight is inversely proportional to its distance from the original pixel's color, raised to the `--power` factor. Colors closer to the original pixel's color receive higher weights.

- **Blended output:** The final color for the pixel is then calculated as a weighted average of these nearest palette colors. This blending process is key to achieving smooth, continuous color transitions.

- **Concurrent processing:** To optimize performance, particularly for large images, the image processing is parallelized. The image is divided into horizontal sections, which are then processed concurrently by separate goroutines.

## Acknowledgments

This project has been inspired by the work of others in the open-source community:

- [Achno/gowall](https://github.com/Achno/gowall)
- [ozwaldorf/lutgen-rs](https://github.com/ozwaldorf/lutgen-rs)
- [nekowinston/faerber](https://github.com/nekowinston/faerber)
- [lighttigerXIV/catppuccinifier](https://github.com/doprz/dipc)
- [doprz/dipc](https://github.com/doprz/dipc)

## Contributing

Contributions are always welcome! If you have ideas, bug reports, or want to submit code, please feel free to open an issue or a pull request.

## License

[MIT License](LICENSE)