package main

import (
	"image/color"
	"machine"
	"time"

	"github.com/aykevl/ledsgo"
	"tinygo.org/x/drivers/ws2812"
)

const NUM_LEDS = 50

var leds = make([]color.RGBA, NUM_LEDS)

const animationSpeed = 2 // higher means faster
var brightness uint8 = 127

type ledPosition struct {
	X uint8
	Y uint8
}

var brightnessMap = [10]uint8{2, 7, 17, 32, 52, 77, 108, 146, 189, 238}

var positions = [...]ledPosition{
	{43, 16}, {44, 16}, {44, 19}, {37, 19}, {39, 9}, {44, 19}, {50, 13}, {41, 12}, {39, 9}, // ball 1
	{57, 20}, {66, 21}, {75, 19}, {64, 15}, {63, 24}, // ball 2
	{65, 28}, {75, 30}, {70, 33}, {80, 40}, {78, 36}, {75, 39}, // ball 3
	{70, 41}, {70, 39}, {62, 37}, {52, 43}, {62, 52}, {66, 40}, {61, 36}, {52, 40}, // ball 4
	{43, 39}, {38, 37}, {45, 36}, {42, 34}, {34, 40}, {42, 44}, {46, 39}, {40, 36}, // ball 5
	{29, 36}, {26, 43}, {29, 36}, // ball 6
	{21, 35}, {15, 37}, {20, 39}, // ball 7
	{29, 29}, {33, 20}, {28, 12}, {23, 20}, {20, 19}, {22, 21}, {28, 27}, {33, 25}, // ball 8
}

func main() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	LED_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	strip := ws2812.New(LED_PIN)

	var command byte
	animation := noise
	for {
		if command == 0 {
			// Read command.
			if machine.Serial.Buffered() != 0 {
				command, _ = machine.Serial.ReadByte()
			}
		}
		if command != 0 {
			switch command {
			case 'N': // noise
				animation = noise
				command = 0
			case 'L': // lightning
				animation = lightning
				command = 0
			case 'D': // disable
				animation = poweroff
				command = 0
			case 'b': // brightness
				if machine.Serial.Buffered() != 0 {
					b, _ := machine.Serial.ReadByte()
					if b >= '0' && b <= '9' {
						brightness = brightnessMap[b-'0']
					}
					command = 0
				}
			}
		}

		// Update colors.
		var t uint64
		if animationSpeed != 0 {
			t = uint64(time.Now().UnixNano() >> (26 - animationSpeed))
		}
		animation(t, leds)

		// Send new colors to LEDs.
		for _, c := range leds {
			strip.WriteByte(c.G) // G
			strip.WriteByte(c.R) // R
			strip.WriteByte(c.B) // B
			strip.WriteByte(0)   // W (alpha channel, used as white channel)
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func noise(t uint64, leds []color.RGBA) {
	for i := range leds {
		pos := positions[i]
		const spread = 24 // higher means more detail
		val := ledsgo.Noise3(uint32(t), uint32(pos.X)*spread, uint32(pos.Y)*spread)
		c := ledsgo.PartyColors.ColorAt(val)
		leds[i] = ledsgo.ApplyAlpha(c, brightness)
	}
}

func lightning(t uint64, leds []color.RGBA) {
	const interval = 1 << 8
	elapsed := interval - t%interval
	for i := range leds[:10] {
		leds[i] = color.RGBA{0, 0, 0, uint8(elapsed / (interval / 100))}
	}
}

func poweroff(t uint64, leds []color.RGBA) {
	for i := range leds {
		leds[i] = color.RGBA{}
	}
}

func xorshift64(x uint64) uint64 {
	// https://en.wikipedia.org/wiki/Xorshift
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17
	return x
}
