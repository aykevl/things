//go:build !baremetal

package main

import (
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl/pixel"
)

func initHardware() {
	board.Simulator.AddressableLEDs = 18
	board.AddressableLEDs.Configure()

	board.Buttons.Configure()
}

func disableLEDs() {
	// Assume LEDs have been shut down already.
	// TODO: the board package should have a way to shut down LEDs when not in
	// use (this is possible on the MCH2022 badge, for example).
	updateLEDs()
}

var simulatorButtonPressed bool

func isButtonPressed() bool {
	for {
		event := board.Buttons.NextEvent()
		if event == board.NoKeyEvent {
			break
		}
		// We assume only one button is used (otherwise this code wouldn't work
		// correctly).
		simulatorButtonPressed = event.Pressed()
	}
	return simulatorButtonPressed
}

func updateLEDs() {
	updateBoardLEDs(board.AddressableLEDs.Data)
	board.AddressableLEDs.Update()
	time.Sleep(time.Second / 500)
}

func updateBoardLEDs[T pixel.Color](data []T) {
	for i, c := range leds {
		data[i] = pixel.NewLinearColor[T](c.R, c.G, c.B)
	}
}
