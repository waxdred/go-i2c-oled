package main

import (
	"fmt"

	"github.com/waxdred/go-i2c-oled/ssd1306"
)

func (i *I2c) SetContrast(contrast int) error {
	var err error
	if contrast < 0 || contrast > 255 {
		return fmt.Errorf("Contrast must be a values from 0 to 255")
	}
	if _, err = i.WriteCommand(OLED_CMD_CONTRAST); err != nil {
		return err
	}
	_, err = i.WriteCommand(byte(contrast))
	i.screen.contrast = contrast
	return err
}

func (i *I2c) SetDim(dim bool) error {
	contrast := i.screen.contrast
	if !dim {
		if i.screen.vccState == ssd1306.SSD1306_EXTERNALVCC {
			contrast = 0x9f
		} else {
			contrast = 0xCF
		}
	}
	return i.SetContrast(contrast)
}
