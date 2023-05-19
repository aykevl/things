package main

import (
	"github.com/aykevl/tinygl/pixel"
)

var leds [18]pixel.LinearGRB888

func main() {
	initLEDs()

	var traceIndex uint8
	var cycle uint8
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
