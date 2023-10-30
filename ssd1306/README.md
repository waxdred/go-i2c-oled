# SSD1306 OLED Driver for Go

This package provides an implementation for controlling an SSD1306 OLED display using Go.

## Features
- Support for different OLED screen dimensions:
    - 128x64
    - 128x32
    - 96x16

- Initialization and basic configuration commands.
- Modular design for easy expansion and integration.

## Usage

1) Initialize the Display

To create a display instance based on your screen's dimensions:
```go
display, err := ssd1306.NewDisplay(width, height, fd, vccstate)
if err != nil {
    log.Fatalf("Failed to create a new display: %v", err)
}
```

2) Initialize the Screen
```go
err = display.Initialize()
if err != nil {
    log.Fatalf("Failed to initialize the screen: %v", err)
}
```

3) You can now proceed with any other SSD1306 related operations.

## Configuration

- VCC State:
    - SSD1306_EXTERNALVCC: External power supply.
    - SSD1306_SWITCHCAPVCC: Internal power supply (default).

- Screen Dimensions:
    - Supported dimensions include: 128x64, 128x32, and 96x16.

## Files
- ssd1306.go: Contains constants and primary functions for controlling the SSD1306.
- SSD1306_96_16.go: Implementation for 96x16 dimension OLED screens.
- SSD1306_128_64.go: Implementation for 128x64 dimension OLED screens.
- SSD1306_128_32.go: Implementation for 128x32 dimension OLED screens.

## Requirements
Go version 1.xx.x (Replace with your version)
Access to the appropriate device interface (typically I2C) on your system.
