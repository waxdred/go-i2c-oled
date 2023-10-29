package ssd1306

import "os"

type SSD1306_96_16 struct {
	fd       *os.File
	vccstate byte
}

// NewSSD1306_96_16 creates a new instance of the SSD1306_96_16 structure.
func NewSSD1306_96_16(fd *os.File, vccstate byte) *SSD1306_96_16 {
	return &SSD1306_96_16{
		fd:       fd,
		vccstate: vccstate,
	}
}

func (d *SSD1306_96_16) Initialize() error {
	data := []byte{
		SSD1306_DISPLAYOFF,         // 0xAE
		SSD1306_SETDISPLAYCLOCKDIV, // 0xD5
		0x60,                       // the suggested ratio 0x60
		SSD1306_SETMULTIPLEX,       // 0xA8
		0x0F,                       // Multiplex value for 96x16
		SSD1306_SETDISPLAYOFFSET,   // 0xD3
		0x0,                        // no offset
		SSD1306_SETSTARTLINE | 0x0, // line #0
		SSD1306_CHARGEPUMP,         // 0x8D
	}

	// Adjust charge pump settings based on vccstate.
	if d.vccstate == SSD1306_EXTERNALVCC {
		data = append(data, byte(0x10)) // External Vcc
	} else {
		data = append(data, byte(0x14)) // Internal Vcc
	}

	// Additional setup commands.
	data = append(data, []byte{
		SSD1306_MEMORYMODE,     // 0x20
		0x00,                   // 0x0 act like ks0108
		SSD1306_SEGREMAP | 0x1, // Map segment 0 to column 95
		SSD1306_COMSCANDEC,     // Scan in descending order
		SSD1306_SETCOMPINS,     // 0xDA
		0x02,                   // Sequential COM pin configuration for 96x16
		SSD1306_SETCONTRAST,    // 0x81
		0x8F,                   // Contrast value
		SSD1306_SETPRECHARGE,   // 0xd9
	}...)

	// Set precharge period based on vccstate.
	if d.vccstate == SSD1306_EXTERNALVCC {
		data = append(data, byte(0x22)) // Precharge value for External Vcc
	} else {
		data = append(data, byte(0xF1)) // Precharge value for Internal Vcc
	}

	// Final setup commands.
	data = append(data, []byte{
		SSD1306_SETVCOMDETECT,       // 0xDB
		0x40,                        // VCOM deselect level
		SSD1306_DISPLAYALLON_RESUME, // 0xA4
		SSD1306_NORMALDISPLAY,       // 0xA6
	}...)

	return sendCommands(d.fd, data...)
}
