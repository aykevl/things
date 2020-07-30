// +build pca10040,v2

// Actually targetting a GT832C_01, which is a small nrf52840 board that's easy
// to solder. The GT832C_01 (as opposed to the GT832E_01) also includes an
// onboard gyroscope/accelerometer (BMI160) which is quite useful.

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
	spiClockPin  machine.Pin = 26
	spiDataPin   machine.Pin = 25
	spiFrequency             = 8000000

	mosfetPin machine.Pin = 27

	serialTxPin machine.Pin = 18

	hasBMI160 = true

	numLeds = 36 // number of LEDs in the strip
	height  = 36 // number of LEDs to be animated
)
