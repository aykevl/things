// +build pca10040,!v2

// Actually targetting a GT832E_01, which is a small nrf52840 board that's easy
// to solder.

package main

import (
	"image/color"
	"machine"
)

// Appearance configuration.
var (
	baseColor     = color.RGBA{0, 0, 0xff, 0x11}
	bluetoothName = "blue poi"
)

// Hardware configuration.
const (
	spiClockPin  machine.Pin = 7
	spiDataPin   machine.Pin = 11
	spiFrequency             = 8000000

	mosfetPin = machine.NoPin

	serialTxPin = machine.NoPin

	hasBMI160 = false

	numLeds = 30 // number of LEDs in the strip
	height  = 14 // number of LEDs to be animated
)
