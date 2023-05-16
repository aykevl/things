package main

import (
	"github.com/aykevl/tinygl/pixel"
)

func main() {
	initLEDs()

	var traceIndex uint8
	var cycle uint8
	animation := uint8(0)
	for {
		cycle++

		switch animation {
		case 0:
			// Trace animation.
			if cycle%64 == 0 {
				traceIndex++
				if traceIndex >= 18 {
					traceIndex = 0
				}
			}
			purpleCircles(leds[:], traceIndex)
		}

		updateLEDs()
	}
}

// Three purple tracers running in circles.
func purpleCircles[T pixel.Color](leds []T, traceIndex uint8) {
	for i := uint8(0); i < 18; i++ {
		idx := i + traceIndex
		if idx >= 18 {
			idx -= 18
		}
		colorIndex := i
		if i >= 6 {
			colorIndex -= 6
		}
		if i >= 12 {
			colorIndex -= 6
		}

		switch colorIndex {
		case 5:
			leds[idx] = makeColor[T](31, 0, 16)
		case 4:
			leds[idx] = makeColor[T](16, 0, 12)
		case 3:
			leds[idx] = makeColor[T](8, 0, 8)
		case 2:
			leds[idx] = makeColor[T](4, 0, 4)
		case 1:
			leds[idx] = makeColor[T](2, 0, 2)
		case 0:
			leds[idx] = makeColor[T](1, 0, 1)
		}
	}
}
