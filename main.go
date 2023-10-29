package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/waxdred/go-i2c-oled/ssd1306"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func main() {
	i2s, err := NewI2c(ssd1306.SSD1306_SWITCHCAPVCC, 64, 128, 0x3C, 1)
	if err != nil {
		fmt.Println("Init err:", err)
		panic(err)
	}
	fmt.Printf("Sreen\nh: %d\nw: %d\noled_char_lenght: %d\noled_max_row: %d\noled_max_cols: %d\n",
		i2s.screen.h, i2s.screen.w, i2s.screen.oled_char_lenght, i2s.screen.oled_max_rows, i2s.oled_max_cols)
	defer i2s.Close()
	// i2s.WriteMsg("Master\nIP: 10.27.27.80\nmem: 12%\ndisk: 1")
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(i2s.img, i2s.img.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)
	colWhite := color.RGBA{255, 255, 255, 255}
	point := fixed.Point26_6{fixed.Int26_6(0 * 64), fixed.Int26_6(15 * 64)} // x = 60, y = 15
	drawer := &font.Drawer{
		Dst:  i2s.img,
		Src:  &image.Uniform{colWhite},
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	drawer.DrawString("Hello from go")
	f, err := os.Create("output.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, i2s.img); err != nil {
		f.Close()
		log.Fatal(err)
	}
	err = i2s.Draw()
	i2s.Clear()
	err = i2s.Display()
	fmt.Println("Display", err)

	// i2s.WriteMsg("Hello Word\n from GO")
	// i2s.WriteChar('t')
	// i2s.WriteChar('e')
}
