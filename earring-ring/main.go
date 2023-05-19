package main

import (
	"github.com/aykevl/ledsgo"
	"github.com/aykevl/tinygl/pixel"
)

var leds [18]pixel.LinearGRB888

var cycle uint16

func main() {
	initLEDs()

	var traceIndex uint8
	animation := uint8(0)
	ledIndex := uint8(0)
	for {
		updateLEDs()

		// Update 3 LEDs.
		for i := uint8(0); i < 3; i++ {
			ledIndex++
			if ledIndex >= 18 {
				ledIndex = 0
				cycle++
				// Trace animation.
				if cycle%4 == 1 {
					traceIndex++
					if traceIndex >= 18 {
						traceIndex = 0
					}
				}
			}

			switch animation {
			case 0:
				purpleCircles(ledIndex, traceIndex)
			case 1:
				rainbowTrace(ledIndex, traceIndex)
			}
		}
	}
}

// Three purple tracers running in circles.
func purpleCircles(index, traceIndex uint8) {
	// This animation has three tracers.
	traceIndex2 := traceIndex + uint8(len(leds))/3
	if traceIndex2 >= uint8(len(leds)) {
		traceIndex2 -= uint8(len(leds))
	}
	traceIndex3 := traceIndex + uint8(len(leds))*2/3
	if traceIndex3 >= uint8(len(leds)) {
		traceIndex3 -= uint8(len(leds))
	}

	if index == traceIndex || index == traceIndex2 || index == traceIndex3 {
		// First tracer.
		leds[index] = pixel.NewLinearGRB888(124, 0, 64)
	} else {
		// dim LED
		c := any(leds[index]).(pixel.LinearGRB888)
		r := uint8(uint16(c.R) * 224 / 256)
		g := uint8(uint16(c.G) * 224 / 256)
		b := uint8(uint16(c.B) * 224 / 256)
		leds[index] = pixel.NewLinearGRB888(r, g, b)
	}
}

func rainbowTrace(index, traceIndex uint8) {
	// This animation has two tracers.
	traceIndex2 := traceIndex + uint8(len(leds))/2
	if traceIndex2 >= uint8(len(leds)) {
		traceIndex2 -= uint8(len(leds))
	}

	if index == traceIndex {
		// First tracer.
		c1 := ledsgo.Color{H: cycle * 128, S: 255, V: 255}.Rainbow()
		leds[index] = pixel.NewLinearGRB888(c1.R, c1.G, c1.B)
	} else if index == traceIndex2 {
		// Second tracer, offset 180° on the color wheel and on the actual LED ring.
		c2 := ledsgo.Color{H: cycle*128 + 0x8000, S: 255, V: 255}.Rainbow()
		leds[traceIndex2] = pixel.NewLinearGRB888(c2.R, c2.G, c2.B)
	} else {
		// dim LED
		c := any(leds[index]).(pixel.LinearGRB888)
		r := uint8(uint16(c.R) * 240 / 256)
		g := uint8(uint16(c.G) * 240 / 256)
		b := uint8(uint16(c.B) * 240 / 256)
		leds[index] = pixel.NewLinearGRB888(r, g, b)
	}
}
