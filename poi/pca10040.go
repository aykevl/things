// +build pca10040

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
)
