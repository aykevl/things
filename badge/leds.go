package main

// Demonstrate support for on-board RGB LEDs.

import (
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/ledsgo"
	"github.com/aykevl/tinygl/pixel"
)

var ledsStop chan struct{}

func init() {
	board.AddressableLEDs.Configure()
}

func toggleLEDs() {
	if ledsStop != nil {
		// Stop running LEDs.
		close(ledsStop)
		ledsStop = nil
		return
	}

	// Start LEDs.
	ledsStop = make(chan struct{})
	go showLEDs(board.AddressableLEDs.Data, ledsStop)
}

func showLEDs[T pixel.Color](data []T, stop chan struct{}) {
	for {
		select {
		case <-stop:
			// An exit was requested.
			// Set all LEDs back to black (off).
			for i := range data {
				data[i] = pixel.NewLinearColor[T](0, 0, 0)
			}
			board.AddressableLEDs.Update()
			return
		default:
			// Continue showing LEDs.
		}

		// Update LEDs.
		now := time.Now()
		for i := range data {
			index := i*4096 + int(now.UnixNano()>>8)
			data[i] = rainbowColor[T](uint16(index))
		}
		board.AddressableLEDs.Update()

		time.Sleep(time.Second / 60)
	}
}

func rainbowColor[T pixel.Color](index uint16) T {
	// Rainbow() returns color.RGBA, but is actually in linear sRGB space.
	// This needs to be fixed eventually.
	c := ledsgo.Color{H: uint16(index), S: 255, V: 255}.Rainbow()
	return pixel.NewLinearColor[T](c.R, c.G, c.B)
}
