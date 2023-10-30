package goi2coled

// Write data to I2C
func (i *I2c) Write(b []byte) (int, error) {
	return i.fd.Write(b)
}

// Send command to OLED
func (i *I2c) WriteCommand(cmd byte) (int, error) {
	return i.Write([]byte{OLED_CMD, cmd})
}

// Send data to OLED
func (i *I2c) WriteData(data []byte) (int, error) {
	res := 0
	for _, value := range data {
		if _, err := i.Write([]byte{OLED_DATA, value}); err != nil {
			return res, err
		}
		res++
	}
	return res, nil
}
