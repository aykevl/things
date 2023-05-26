package main

import (
	"image/color"
	"machine"
	"runtime/interrupt"
	"time"

	"github.com/aykevl/ledsgo"
	"tinygo.org/x/drivers/ws2812"
)

const NUM_LEDS = 10

var leds = make([]color.RGBA, NUM_LEDS)

func main() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	LED_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	strip := ws2812.New(LED_PIN)

	for {
		// Update colors.
		t := uint64(time.Now().UnixNano() >> 20)
		noise(t, leds)

		// Send new colors to LEDs.
		mask := interrupt.Disable()
		for _, c := range leds {
			strip.WriteByte(c.G) // G
			strip.WriteByte(c.R) // R
			strip.WriteByte(c.B) // B
			strip.WriteByte(c.A) // W (alpha channel, used as white channel)
		}
		interrupt.Restore(mask)
		time.Sleep(20 * time.Millisecond)
	}
}

func noise(t uint64, leds []color.RGBA) {
	for i := range leds {
		const spread = 48 // higher means more detail
		val := ledsgo.Noise2(uint32(t), uint32(i)*spread)
		c := ledsgo.PartyColors.ColorAt(val)
		c.A = 0 // the alpha channel is used as white channel, so don't use it
		leds[i] = ledsgo.ApplyAlpha(c, 64)
	}
}
