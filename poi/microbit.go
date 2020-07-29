// +build microbit

// Actually targetting a GT82C_02, which is a small nrf51822 board that's easy
// to solder. Peripherals are kept at the same pins whenever possible.

package main

import (
	"image/color"
	"machine"
)

// Appearance configuration.
var (
	baseColor     = color.RGBA{0xff, 0, 0, 0x11}
	bluetoothName = "red poi"
)

// Hardware configuration.
const (
	spiFrequency = 8000000
	spiClockPin  = machine.SPI0_SCK_PIN
	spiDataPin   = machine.SPI0_SDO_PIN
)
