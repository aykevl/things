package main

import (
	"image/color"
	"machine"
	"time"

	"github.com/aykevl/ledsgo"
	"tinygo.org/x/drivers/ws2812"
)

const NumLEDs = 8 + 8

const (
	ledPin = machine.GPIO28
)

var leds = make([]color.RGBA, NumLEDs)

var brightness uint8 = 24

var palette ledsgo.Palette16 = ledsgo.RainbowColors

type Point struct {
	x, y int8
}

// Map 8 LEDs in a circle to (x, y) coordinates.
var mapping = [...]Point{
	{127, 0},
	{90, 90},
	{0, 127},
	{-90, 90},
	{-127, 0},
	{-90, -90},
	{0, -127},
	{90, -90},
}

func main() {
	ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	strip := ws2812.New(ledPin)

	for {
		// Update colors.
		now := time.Now()
		t := uint64(now.UnixNano() >> 22)
		noise(t, leds)

		// Send new colors to LEDs.
		strip.WriteColors(leds)
		time.Sleep(10 * time.Millisecond)
	}
}

func noise(t uint64, leds []color.RGBA) {
	for i := range leds {
		x := (int(mapping[i%8].x) + 127) * 4
		y := (int(mapping[i%8].y) + 127) * 4
		val := ledsgo.Noise3(uint32(x)+(uint32(i)/8)<<16, uint32(y), uint32(t))
		c := palette.ColorAt(val)
		c.A = 0 // the alpha channel is used as white channel, so don't use it
		leds[i] = ledsgo.ApplyAlpha(c, brightness)
	}
}
