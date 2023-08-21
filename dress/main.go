package main

import (
	"image/color"
	"machine"
	"time"

	"github.com/aykevl/ledsgo"
	"tinygo.org/x/drivers/ws2812"
)

const NumLEDs = 50

const (
	//ledPin = machine.ADC5 // arduino
	ledPin = machine.GPIO29
)

var leds = make([]color.RGBA, NumLEDs)

var brightness uint8 = 32

// var palette = ledsgo.PartyColors
var palette = PurpleShades

var PurpleShades = ledsgo.Palette16{
	color.RGBA{0xFF, 0x00, 0x00, 0xFF},
	color.RGBA{0xFF, 0x00, 0x00, 0xFF},
	color.RGBA{0xEE, 0x00, 0x11, 0xFF},
	color.RGBA{0xDD, 0x00, 0x22, 0xFF},
	color.RGBA{0xBB, 0x00, 0x44, 0xFF},
	color.RGBA{0x88, 0x00, 0x66, 0xFF},
	color.RGBA{0x77, 0x00, 0x77, 0xFF},
	color.RGBA{0x77, 0x00, 0x66, 0xFF},
	color.RGBA{0x66, 0x00, 0x66, 0xFF},
	color.RGBA{0x66, 0x00, 0x77, 0xFF},
	color.RGBA{0x55, 0x00, 0x88, 0xFF},
	color.RGBA{0x55, 0x00, 0x88, 0xFF},
	color.RGBA{0x44, 0x00, 0x99, 0xFF},
	color.RGBA{0x33, 0x00, 0xAA, 0xFF},
	color.RGBA{0x22, 0x00, 0xBB, 0xFF},
	color.RGBA{0x11, 0x00, 0xEE, 0xFF},
}

func main() {
	ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	strip := ws2812.New(ledPin)

	for {
		// Update colors.
		t := uint64(time.Now().UnixNano() >> 20)
		noise(t, leds)
		//showPalette(t, leds)

		// Send new colors to LEDs.
		strip.WriteColors(leds)
		time.Sleep(1 * time.Millisecond)
	}
}

func position(i int) (x, y int) {
	x = i / 10 * 2
	loopPos := i - (x * 5)
	y = loopPos
	if y >= 5 {
		y = 9 - y
		x += 1
	}
	return
}

func noise(t uint64, leds []color.RGBA) {
	for i := range leds {
		const spread = 7 // higher means more detail
		x, y := position(i)
		val := ledsgo.Noise3(uint32(x)<<spread, uint32(t), uint32(t)-(uint32(y)<<spread)+0x1000)
		c := palette.ColorAt(val)
		c.A = 0 // the alpha channel is used as white channel, so don't use it
		leds[i] = ledsgo.ApplyAlpha(c, brightness)
	}
}

func showPalette(t uint64, leds []color.RGBA) {
	for i := range leds {
		const spread = 7 // higher means more detail
		val := uint16(i * 256 / NumLEDs)
		c := palette.ColorAt(val * 256)
		c.A = 0 // the alpha channel is used as white channel, so don't use it
		leds[i] = ledsgo.ApplyAlpha(c, brightness)
	}
}
