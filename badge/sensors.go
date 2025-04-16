package main

import (
	"strconv"
	"time"
	"runtime"

	"github.com/aykevl/board"
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/style"
	"github.com/aykevl/tinygl/style/basic"
	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/pixel"
)

func showSensors[T pixel.Color](screen *tinygl.Screen[T]) {
	// Determine size and scale of the screen.
	scalePercent := board.Display.PPI() * 100 / 120

	// Create UI.
	theme := basic.New(style.NewScale(scalePercent), screen)
	header := theme.NewText("Sensors")
	header.SetBackground(pixel.NewColor[T](0, 0, 255))
	header.SetColor(pixel.NewColor[T](255, 255, 255))
	battery := theme.NewText("battery: ...")
	battery.SetAlign(tinygl.AlignLeft)
	voltage := theme.NewText("voltage: ...")
	voltage.SetAlign(tinygl.AlignLeft)
	temperature := theme.NewText("temp: ...")
	temperature.SetAlign(tinygl.AlignLeft)
	acceleration := theme.NewText("accel: ...")
	acceleration.SetAlign(tinygl.AlignLeft)
	steps := theme.NewText("steps: ...")
	steps.SetAlign(tinygl.AlignLeft)
	cpus := theme.NewText("CPU cores: " + strconv.Itoa(runtime.NumCPU()))
	cpus.SetAlign(tinygl.AlignLeft)
	vbox := theme.NewVBox(header, battery, voltage, temperature, acceleration, steps, cpus)
	screen.SetChild(vbox)

	// Show screen.
	screen.Update()

	// Initialize sensors.
	board.Power.Configure()
	board.Sensors.Configure(drivers.Acceleration | drivers.Temperature)

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

		// Read battery status.
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

		// Read sensors.
		board.Sensors.Update(drivers.Temperature | drivers.Acceleration)
		celsius := (board.Sensors.Temperature() + 500) / 1000
		temperature.SetText("temp: " + strconv.Itoa(int(celsius)) + "C")
		ax, ay, az := board.Sensors.Acceleration()
		acceleration.SetText("accel: " + strconv.FormatInt(int64(ax/10000), 10) + " " + strconv.FormatInt(int64(ay/10000), 10) + " " + strconv.FormatInt(int64(az/10000), 10))
		steps.SetText("steps: " + strconv.Itoa(int(board.Sensors.Steps())))

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
