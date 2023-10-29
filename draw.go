package main

import "fmt"

func (i *I2c) convertImageToOLEDData() ([]byte, error) {
	size := i.screen.w * i.screen.h / PIXSIZE
	data := make([]byte, size)
	for y := 0; y < int(i.screen.h); y++ {
		for x := 0; x < int(i.screen.w); x++ {
			pixel := i.img.At(x, y)
			r, g, b, _ := pixel.RGBA()
			avg := (r + g + b) / 3
			if avg > 0x7FFF {
				byteIndex := y*int(i.screen.w)/8 + x/8
				bitIndex := x % 8
				data[byteIndex] |= 1 << bitIndex
			}
		}
	}
	fmt.Println(data)
	return data, nil
}

// Drawn image to Oled
func (i *I2c) Draw() error {
	var err error
	i.buffer, err = i.convertImageToOLEDData()
	if err != nil {
		return err
	}
	return nil
}
