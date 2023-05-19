//go:build !baremetal

package main

import (
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl/pixel"
)

var leds []pixel.RGB888

func initLEDs() {
	board.Simulator.AddressableLEDs = 18
	board.AddressableLEDs.Configure()
	leds = board.AddressableLEDs.Data
}

func updateLEDs() {
	board.AddressableLEDs.Update()
	time.Sleep(time.Second / 2400)
}

func makeColor[T pixel.Color](r, g, b uint8) T {
	return pixel.NewLinearColor[T](r, g, b)
}
