# goi2coled: A Go library for OLED screens based on SSD1306 via I2C

`goi2coled` is a Go library designed to facilitate interactions with OLED screens that use the SSD1306 controller via the I2C bus. With the integration of Go's image package, this library allows developers to design and display graphics on OLED screens with ease.

## Features
- Initialisation of I2C connections to the OLED screen.
- Displaying and clearing content on the screen.
- Adjusting the display contrast and dimming.
- Directly drawing images from Go's image.RGBA to the OLED screen.
- Convert images to data compatible with the OLED screen.
- Direct communication with the SSD1306 controller using commands and data writes.

## Prerequisites
An environment setup to run Go code.
The SSD1306 OLED screen and compatible I2C connections.
Dependency: github.com/waxdred/go-i2c-oled/ssd1306

## Usage
Here's a brief guide on how to use goi2coled:

1) Initialising the I2C connection:
```go
i2c, err := NewI2c(vccState, height, width, address, bus)
if err != nil {
    log.Fatal(err)
}
defer i2c.Close()
```

2) Drawing an Image:
```go
img := image.NewRGBA(image.Rect(0, 0, width, height))
// ... (draw or manipulate the img as needed) ...
i2c.DrawImage(img)
```

3) Displaying Buffer to Screen:
```go
err = i2c.Display()
if err != nil {
    log.Fatal(err)
}
```

4) Adjusting Contrast:
```go
err = i2c.SetContrast(contrastValue)  // contrastValue: 0 to 255
```

5) Dimming the Display:
```go
err = i2c.SetDim(true)  // Set to true to dim
```

## Exemple
```go
import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

        "github.com/waxdred/go-i2c-oled"
	"github.com/waxdred/go-i2c-oled/ssd1306"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func main() {
	// Initialize the OLED display with the provided parameters
	oled, err := goi2coled.NewI2c(ssd1306.SSD1306_SWITCHCAPVCC, 64, 128, 0x3C, 1)
	if err != nil {
		panic(err)
	}
	defer oled.Close()


    	// Ensure the OLED is properly closed at the end of the program
	defer oled.Close()

    	// Define a black color
	black := color.RGBA{0, 0, 0, 255}

    	// Set the entire OLED image to black
	draw.Draw(oled.Img, oled.Img.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)

    	// Define a white color
	colWhite := color.RGBA{255, 255, 255, 255}

    	// Set the starting point for drawing text
	point := fixed.Point26_6{fixed.Int26_6(0 * 64), fixed.Int26_6(15 * 64)} // x = 0, y = 15

    	// Configure the font drawer with the chosen font and color
	drawer := &font.Drawer{
		Dst:  oled.Img,
		Src:  &image.Uniform{colWhite},
		Face: basicfont.Face7x13,
		Dot:  point,
	}

    	// Clear the OLED image (making it all black)
	draw.Draw(oled.Img, oled.img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

    	// Draw the text "Hello" on the OLED image
	drawer.DrawString("Hello")
    
    	// Move the drawing point down by 10 pixels for the next line of text
	drawer.Dot.Y += fixed.Int26_6(10 * 64)

    	// Set the drawing point's x coordinate back to 0 for alignment
	drawer.Dot.X = fixed.Int26_6(0 * 64)

    	// Draw the text "From golang!" on the OLED image
	drawer.DrawString("From golang!")

    	// Clear the OLED's buffer (if applicable to your library)
	oled.Clear()

    	// Update the OLED's buffer with the current image data
	oled.Draw()

    	// Display the buffered content on the OLED screen
	err = oled.Display()
}
```

## Notes
- Ensure your SSD1306 OLED is properly connected via I2C.
- Be aware of the screen resolution; the provided library assumes specific width and height, which should match your OLED screen.

## Contributions
Contributions, bug reports, and feature requests are welcome! Feel free to open an issue or submit a pull request.

