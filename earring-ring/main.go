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
	for {
		cycle++

		switch animation {
		case 0:
			// Trace animation.
			if cycle%16 == 0 {
				traceIndex++
				if traceIndex >= 18 {
					traceIndex = 0
				}
			}
			purpleCircles(traceIndex)
		}

		updateLEDs()
	}
}

// Three purple tracers running in circles.
func purpleCircles(traceIndex uint8) {
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
			leds[idx] = pixel.NewLinearGRB888(124, 0, 64)
		case 4:
			leds[idx] = pixel.NewLinearGRB888(64, 0, 48)
		case 3:
			leds[idx] = pixel.NewLinearGRB888(32, 0, 32)
		case 2:
			leds[idx] = pixel.NewLinearGRB888(16, 0, 16)
		case 1:
			leds[idx] = pixel.NewLinearGRB888(8, 0, 8)
		case 0:
			leds[idx] = pixel.NewLinearGRB888(4, 0, 4)
		}
	}
}
