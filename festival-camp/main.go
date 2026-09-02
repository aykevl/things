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
	ledPin    = machine.GPIO29
	buttonPin = machine.GPIO24
)

var leds = make([]color.RGBA, NumLEDs)

var brightness uint8 = 63

var paletteIndex uint8
var palettes = [...]ledsgo.Palette16{
	ledsgo.PartyColors,
	ledsgo.RainbowColors,
	ledsgo.ForestColors,
	ledsgo.OceanColors,
	ledsgo.CloudColors,
}

func main() {
	ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	strip := ws2812.New(ledPin)

	buttonPin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	pressedTimes := 0 // TODO: make this a timestamp or something
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

		// Check for button press
		const paletteChangeTimeout = 200
		pressed := !buttonPin.Get()
		if pressedTimes == paletteChangeTimeout {
			paletteIndex++
			if int(paletteIndex) >= len(palettes) {
				paletteIndex = 0
			}
			println("palette", paletteIndex)
		} else if !pressed && pressedTimes > 0 && pressedTimes < paletteChangeTimeout {
			// Next brightness level
			if brightness <= 8 {
				brightness = 255
			} else {
				brightness /= 2
			}
			println("pressed", brightness)
		}
		if pressed {
			pressedTimes++
		} else {
			pressedTimes = 0
		}
	}
}

func noise(t uint64, leds []color.RGBA) {
	palette := &palettes[paletteIndex]
	for i := range leds {
		const spread = 8 // higher means more detail
		val := ledsgo.Noise2(uint32(i)<<spread, uint32(t))
		c := palette.ColorAt(val)
		c.A = 0 // the alpha channel is used as white channel, so don't use it
		leds[i] = ledsgo.ApplyAlpha(c, brightness)
	}
}
