package main

import (
	"image/color"
	"machine"
	"time"

	"github.com/aykevl/ledsgo"
	"tinygo.org/x/drivers/ws2812"
)

const NumLEDs = 30

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

		// Send new colors to LEDs.
		strip.WriteColors(leds)
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
