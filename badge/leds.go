package main

// Demonstrate support for on-board RGB LEDs.

import (
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/ledsgo"
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
	go showLEDs(ledsStop)
}

func showLEDs(stop chan struct{}) {
	array := board.AddressableLEDs
	for {
		select {
		case <-stop:
			// An exit was requested.
			// Set all LEDs back to black (off).
			for i := 0; i < array.Len(); i++ {
				array.SetRGB(i, 0, 0, 0)
			}
			array.Update()
			return
		default:
			// Continue showing LEDs.
		}

		// Update LEDs.
		now := time.Now()
		for i := 0; i < array.Len(); i++ {
			index := i*4096 + int(now.UnixNano()>>8)
			c := ledsgo.Color{H: uint16(index), S: 255, V: 255}.Rainbow()
			array.SetRGB(i, c.R, c.G, c.B)
		}
		array.Update()

		time.Sleep(time.Second / 60)
	}
}
