package main

import (
	"strconv"
	"time"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/pixel"
	"github.com/aykevl/tinygl/style"
	"github.com/aykevl/tinygl/style/basic"
)

func showSensors[T pixel.Color](screen *tinygl.Screen[T]) {
	// Determine size and scale of the screen.
	scalePercent := board.Display.PPI() * 100 / 120

	// Create UI.
	theme := basic.NewTheme(style.NewScale(scalePercent), screen)
	header := theme.NewText("Sensors")
	header.SetBackground(pixel.NewColor[T](0, 0, 255))
	header.SetColor(pixel.NewColor[T](255, 255, 255))
	battery := theme.NewText("battery: ...")
	battery.SetAlign(tinygl.AlignLeft)
	voltage := theme.NewText("voltage: ...")
	voltage.SetAlign(tinygl.AlignLeft)
	vbox := theme.NewVBox(header, battery, voltage)
	screen.SetChild(vbox)

	// Show screen.
	screen.Update()

	// Initialize sensors.
	board.Power.Configure()

	for {
		// Read keyboard inputs.
		board.Buttons.ReadInput()
		for {
			// Read keyboard events.
			event := board.Buttons.NextEvent()
			if event == board.NoKeyEvent {
				break
			}
			if event.Pressed() {
				switch event.Key() {
				case board.KeyB, board.KeyEscape:
					return
				}
			}
		}

		// Read sensors.
		state, microvolts, percent := board.Power.Status()
		batteryText := "battery: "
		if percent >= 0 {
			batteryText += strconv.Itoa(int(percent)) + "% "
		}
		if state != board.Discharging {
			batteryText += state.String()
		}
		battery.SetText(batteryText)
		voltage.SetText("voltage: " + formatVoltage(microvolts))

		screen.Update()
		time.Sleep(time.Second / 5)
	}
}

func formatVoltage(microvolts uint32) string {
	volts := strconv.Itoa(int(microvolts / 1000_000))
	decimals := strconv.Itoa(int(microvolts % 1000_000 / 10_000))
	for len(decimals) < 2 {
		decimals = "0" + decimals
	}
	return volts + "." + decimals + "V"
}
