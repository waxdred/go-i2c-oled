package goi2coled

// Turn off OLED display
func (i *I2c) DisplayOff() (int, error) {
	return i.WriteCommand(OLED_CMD_DISPLAY_OFF)
}

// Turn on OLED display
func (i *I2c) DisplayOn() (int, error) {
	return i.WriteCommand(OLED_CMD_DISPLAY_ON)
}

// Display buffer to the screen
func (i2c *I2c) Display() error {
	i2c.WriteCommand(OLED_CMD_COL_ADDRESSING)
	i2c.WriteCommand(0)
	i2c.WriteCommand(byte(i2c.screen.w - 1))
	i2c.WriteCommand(OLED_CMD_PAGE_ADDRESSING)
	i2c.WriteCommand(0)
	i2c.WriteCommand(byte(i2c.screen.h / 8))
	for i := 0; i < len(i2c.buffer); i += 16 {
		data := i2c.buffer[i : i+16]
		_, err := i2c.WriteData(data)
		if err != nil {
			return err
		}
	}
	return nil
}
