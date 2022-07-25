package main

import (
	"image/color"
	"machine"
	"time"

	"github.com/aykevl/ledsgo"
	"tinygo.org/x/drivers/ws2812"
)

func main() {
	// Enable power to the LEDs
	power := machine.PowerOn
	power.Configure(machine.PinConfig{Mode: machine.PinOutput})
	power.High()

	// Initialize the data pin.
	led := machine.WS2812
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Prepare the LED data slice to send.
	ws := ws2812.New(led)
	leds := make([]color.RGBA, 5)

	// Update the data each cycle.
	for {
		now := time.Now()
		for i := range leds {
			leds[i] = ledsgo.Color{
				H: uint16(i)*4096 + uint16(now.UnixNano()>>17),
				S: 255,
				V: 64,
			}.Rainbow()
		}
		ws.WriteColors(leds)
		time.Sleep(time.Second / 60)
	}
}
