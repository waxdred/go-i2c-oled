package main

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
	prefixedData := make([]byte, 2*len(data))
	for idx, value := range data {
		prefixedData[2*idx] = OLED_DATA
		prefixedData[2*idx+1] = value
	}
	return i.Write(prefixedData)
}
