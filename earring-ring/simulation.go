//go:build !baremetal

package main

import (
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl/pixel"
)

func initLEDs() {
	board.Simulator.AddressableLEDs = 18
	board.AddressableLEDs.Configure()
}

func updateLEDs() {
	updateBoardLEDs(board.AddressableLEDs.Data)
	board.AddressableLEDs.Update()
	time.Sleep(time.Second / 400)
}

func updateBoardLEDs[T pixel.Color](data []T) {
	for i, c := range leds {
		data[i] = pixel.NewLinearColor[T](c.R, c.G, c.B)
	}
}
