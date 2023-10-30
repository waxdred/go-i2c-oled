package goi2coled

import (
	"fmt"
	"image"
	"image/color"
)

func (i *I2c) convertImageToOLEDData() ([]byte, error) {
	bounds := i.Img.Bounds()
	if bounds.Max.X != i.screen.w || i.screen.h != bounds.Max.Y {
		panic(fmt.Sprintf("Error: Size of image is not %dx%d pixels.", i.screen.w, i.screen.h))
	}
	size := i.screen.w * i.screen.h / PIXSIZE
	data := make([]byte, size)
	for page := 0; page < i.screen.h/8; page++ {
		for x := 0; x < i.screen.w; x++ {
			bits := uint8(0)
			for bit := 0; bit < 8; bit++ {
				y := page*8 + 7 - bit
				if y < i.screen.h {
					col := color.GrayModel.Convert(i.Img.At(x, y)).(color.Gray)
					if col.Y > 127 {
						bits = (bits << 1) | 1
					} else {
						bits = bits << 1
					}
				}
			}
			index := page*i.screen.w + x
			data[index] = byte(bits)
		}
	}
	return data, nil
}

func (i *I2c) DrawImage(image *image.RGBA) {
	i.Img = image
	i.buffer, _ = i.convertImageToOLEDData()
}

func (i *I2c) Draw() {
	i.buffer, _ = i.convertImageToOLEDData()
}
