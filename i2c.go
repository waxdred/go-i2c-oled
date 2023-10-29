package main

import (
	"fmt"
	"image"
	"os"
	"syscall"

	"github.com/waxdred/go-i2c-oled/ssd1306"
	// "golang.org/x/image/font"
)

// Constants for OLED commands and addressing const (
const (
	I2C_SLAVE = 0x0703

	OLED_CMD                 = 0x80
	OLED_CMD_COL_ADDRESSING  = 0x21
	OLED_CMD_PAGE_ADDRESSING = 0x22
	OLED_CMD_CONTRAST        = 0x81
	OLED_CMD_START_COLUMN    = 0x00
	OLED_CMD_DISPLAY_OFF     = 0xAE
	OLED_CMD_DISPLAY_ON      = 0xAF

	OLED_DATA            = 0x40
	OLED_ADRESSING       = 0x21
	OLED_ADRESSING_START = 0xB0
	OLED_ADRESSING_COL   = 0x21
	OLED_END             = 0x10
	PIXSIZE              = 8
)

// Struct for representing screen properties
type screen struct {
	h                int
	w                int
	oled_char_lenght int
	oled_max_rows    int
	oled_max_cols    int
	contrast         int
	buffer           []byte
	vccState         int
	img              *image.RGBA
}

// Function to initialize a new screen with given height and width
func newScreen(vccState, h, w int) screen {
	// TODO: Create function to derive oled_char_lenght and oled_max_rows from h and w
	oled_max_rows := int(h / 8)
	oled_max_cols := int(w / 6)
	oled_char_lenght := int(w / oled_max_cols)
	return screen{
		h:                h,
		w:                w,
		oled_max_rows:    oled_max_rows,
		oled_char_lenght: oled_max_cols,
		oled_max_cols:    oled_char_lenght,
		img:              image.NewRGBA((image.Rect(0, 0, int(w), int(h)))),
		vccState:         vccState,
	}
}

// Struct for managing I2C operations
type I2c struct {
	address    int
	bus        int
	fd         *os.File
	currentRow byte
	currentCol byte
	screen
}

// Function to initialize I2C with given parameters
func NewI2c(vccState, h, w, address, bus int) (*I2c, error) {
	fd, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, os.ModeExclusive)
	if err != nil {
		return nil, err
	}
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, fd.Fd(), I2C_SLAVE, uintptr(address), 0, 0, 0)
	if errno != 0 {
		err = syscall.Errno(errno)
		fd.Close()
		return nil, err
	}
	display, err := ssd1306.NewDisplay(int(w), int(h), fd, byte(vccState))
	if err != nil {
		return nil, err
	}
	display.Initialize()

	i2c := &I2c{
		address: address,
		bus:     bus,
		fd:      fd,
		screen:  newScreen(vccState, h, w),
	}
	i2c.DisplayOn()
	return i2c, nil
}

// Close the I2C connection
func (i *I2c) Close() error {
	return i.fd.Close()
}

// Read data from I2C
func (i *I2c) Read(b []byte) (int, error) {
	return i.fd.Read(b)
}

// Set the cursor position on OLED
func (i *I2c) SetCursor(row int, column int) error {
	_, err := i.WriteCommand(byte(OLED_ADRESSING_START + row))
	if err != nil {
		return err
	}
	_, err = i.WriteCommand(byte(OLED_CMD_START_COLUMN + (int(i.screen.oled_char_lenght) * column & 0x0F)))
	if err != nil {
		return err
	}
	_, err = i.WriteCommand(byte(OLED_END + ((int(i.screen.oled_char_lenght) * column >> 4) & 0x0F)))
	if err != nil {
		return err
	}
	i.currentRow = byte(row)
	return err
}

// Set column addressing on OLED
func (i *I2c) SetColumnAddressing(startPX, endPX int) (int, error) {
	var res int

	if res, err := i.WriteCommand(OLED_ADRESSING_COL); err != nil {
		return res, err
	}
	if res, err := i.WriteCommand(byte(startPX)); err != nil {
		return res, err
	}
	if res, err := i.WriteCommand(byte(endPX)); err != nil {
		return res, err
	}
	i.currentRow = 0
	return res, nil
}

// Write empty characters to OLED
func (i *I2c) writeEmpyChars() error {
	for row := 0; row < int(i.screen.oled_max_rows); row++ {
		if err := i.SetCursor(row, 0); err != nil {
			return err
		}
		for px := 0; px < int(i.screen.w); px++ {
			if _, err := i.WriteData([]byte{0x00}); err != nil {
				return err
			}
		}
	}
	i.currentRow = 0
	i.currentCol = 0
	return nil
}

// Clear the OLED screen
func (i *I2c) Clear() {
	size := i.screen.w * i.screen.h / PIXSIZE
	i.buffer = make([]byte, size)
}

// Write a single character to the OLED screen.
// This function looks up the character's pattern in the Table_ascii and sends it to the screen.
func (i *I2c) WriteChar(c int) error {
	index := c - 32
	if index < 0 || index >= len(Table_ascii) {
		return fmt.Errorf("Character out of bounds for Table_ascii")
	}
	dataToWrite := make([]byte, len(Table_ascii[index]))
	for idx, t := range Table_ascii[index] {
		dataToWrite[idx] = byte(t)
	}
	_, err := i.WriteData(dataToWrite)
	if err != nil {
		return err
	}

	i.currentCol++
	if i.currentCol > byte(i.screen.oled_max_cols) {
		i.currentRow++
		i.currentCol = 0
	}
	return nil
}

// Write a message string to the OLED screen.
// This function handles newline characters and writes each character using WriteChar.
func (i2c *I2c) WriteMsg(msg string) (int, error) {
	var res int
	var totalData []byte

	for _, m := range msg {
		if m == '\n' {
			i2c.currentRow++
			if err := i2c.SetCursor(int(i2c.currentRow), 0); err != nil {
				return res, err
			}
			continue
		}
		index := int(m) - 32
		if index < 0 || index >= len(Table_ascii) {
			return res, fmt.Errorf("Character out of bounds for Table_ascii")
		}
		for _, t := range Table_ascii[index] {
			totalData = append(totalData, byte(t))
		}
		i2c.currentCol++
		if i2c.currentCol > byte(i2c.screen.oled_max_cols) {
			i2c.currentRow++
			i2c.currentCol = 0
		}
		res++
	}
	_, err := i2c.WriteData(totalData)
	if err != nil {
		return res, err
	}
	return res, nil
}
