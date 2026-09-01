package main

import (
	"image/color"
	"machine"
	"runtime/interrupt"
	"time"

	"github.com/aykevl/ledsgo"
	"tinygo.org/x/drivers/ws2812"
)

const NumLEDs = 50

const (
	ledPin = machine.GPIO29
)

var leds = make([]color.RGBA, NumLEDs)

var brightness uint8 = 64

var palette = ledsgo.PartyColors

func main() {
	ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	strip := ws2812.New(ledPin)

	for {
		// Update colors.
		t := uint64(time.Now().UnixNano() >> 24)
		noise(t, leds)

		// Test pattern
		//leds[0] = color.RGBA{R: 255}
		//leds[1] = color.RGBA{G: 255}
		//leds[2] = color.RGBA{B: 255}

		// Send new colors to LEDs.
		mask := interrupt.Disable()
		for _, c := range leds[:] {
			strip.WriteByte(c.R)
			strip.WriteByte(c.G)
			strip.WriteByte(c.B)
		}
		interrupt.Restore(mask)
		time.Sleep(1 * time.Millisecond)
	}
}

func noise(t uint64, leds []color.RGBA) {
	for i := range leds {
		const spread = 8 // higher means more detail
		val := ledsgo.Noise2(uint32(i)<<spread, uint32(t))
		c := palette.ColorAt(val)
		c.A = 0 // the alpha channel is used as white channel, so don't use it
		leds[i] = ledsgo.ApplyAlpha(c, brightness)
	}
}
