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
	spiClockPin  = machine.SPI0_SCK_PIN
	spiDataPin   = machine.SPI0_SDO_PIN
	spiFrequency = 8000000

	mosfetPin = machine.NoPin

	serialTxPin = machine.NoPin

	hasBMI160 = false

	numLeds = 30 // number of LEDs in the strip
	height  = 14 // number of LEDs to be animated
)

//go:inline
func setLED(y int16, c color.RGBA) {
	// This is a poi with two sides that wraps around.
	// Make sure the other side is also colored properly.
	leds[y] = c
	leds[len(leds)-int(y)-1] = c
}
