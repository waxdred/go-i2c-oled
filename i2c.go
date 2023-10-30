package goi2coled

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"syscall"

	"github.com/waxdred/go-i2c-oled/ssd1306"
)

// Constants for OLED commands and addressing const (
const (
	I2C_SLAVE = 0x0703

	OLED_CMD                 = 0x80
	OLED_CMD_COL_ADDRESSING  = 0x21
	OLED_CMD_PAGE_ADDRESSING = 0x22
	OLED_CMD_CONTRAST        = 0x81
	OLED_CMD_START_COLUMN    = 0x00
	OLED_CMD_HIGH_COLUMN     = 0x10
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
	h        int
	w        int
	contrast int
	buffer   []byte
	vccState int
	img      draw.Image
}

// Function to initialize a new screen with given height and width
func newScreen(vccState, h, w int) screen {
	return screen{
		h:        h,
		w:        w,
		img:      image.NewRGBA((image.Rect(0, 0, int(w), int(h)))),
		vccState: vccState,
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

// Clear the OLED screen
func (i *I2c) Clear() {
	size := i.screen.w * i.screen.h / PIXSIZE
	i.buffer = make([]byte, size)
}
