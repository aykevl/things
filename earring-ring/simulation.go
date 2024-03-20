//go:build !baremetal

package main

import (
	"time"

	"github.com/aykevl/board"
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
	for i, c := range leds {
		board.AddressableLEDs.SetRGB(i, c.R, c.G, c.B)
	}
	board.AddressableLEDs.Update()
	time.Sleep(time.Second / 500)
}
